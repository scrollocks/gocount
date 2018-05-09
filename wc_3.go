/*

This program counts words in text files, given as command line parameters.

There is concurrency using goroutines, and a channel for the results of word counting.

Author: Jamie Skipworth

*/
package main

import (
	"fmt"
	"io"
	"os"
	"unicode"
	"time"
	"sync"
	)

type Word struct {
	fileName string
	words int64
}

const chunkSize int32 = 32*1024
var filename string

// Writes text to a file descriptor. Nice for consistent logging.
func logLine( fd *os.File, message string ) {
	fmt.Fprintf( fd, "%s - %s\n", time.Now().Format("20060106 15:04:05.000"), message )
}

// This function opens a text file, reads the data in chunks and extracts words from it.
// It sends the word to a results channel.
func wordCounter( filename string, waitGroup *sync.WaitGroup, out chan <- Word ) {

	defer waitGroup.Done()

	// Call Done() when this function ends
	timeStart := time.Now()

	logLine( os.Stderr, fmt.Sprintf( "wordReader started for %s", filename ) )

	// Attempt to open the file
	file, err := os.Open( filename )
	defer file.Close()

	if( err != nil ){
		logLine( os.Stderr, fmt.Sprintf( "Error opening file %s", filename ) )
		return
	}

	var totalBytesRead, wordCount int64	// Bytes read
	var str, lastWord string	// str is an individual character. lastWord is a word.

	// Start processing data from the file
	for {

		buf := make( []byte, chunkSize )	// Create a buffer
		bytes, err := file.Read( buf )		// Read data into the buffer

		// Count words only if we read some data.
		if bytes > 0 {
			totalBytesRead += int64(bytes)	// Sum the bytes we've read 
			str = string( buf )				// Convert bytes to a string

			// Boolean that tells us if we're in a region of white space.
			var inSpace bool = true

			// Convert the string into an array of runes (unicode) and iterate through it a char at a time
			for _, r := range []rune( str )  {

				// Is the rune white-space? Set the isSpace flag if it is - we've entered a region of white space.
				if unicode.IsSpace( r ){
					// If we've just entered a region of white-space and we weren't before, then we must have 
					// encountered a word.
					if ! inSpace {	
						inSpace = true
						wordCount += 1
						lastWord = ""
					}
				// We're in a region of text
				}else{
					inSpace = false
					lastWord += string( r )
				}
			}
		}

		// Did we hit the end of file?
		if err == io.EOF {
			break
		}

		// Did we hit a different error?
		if err != nil {
			logLine( os.Stderr, fmt.Sprintf( "Error reading file %s: %s",  filename, err ) )
			break
		}

	}

	out <- Word{ file.Name(), wordCount }
	file.Close()
	timeEnd := time.Now()
	logLine( os.Stderr, fmt.Sprintf( "wordCounter finished in %s for %s", timeEnd.Sub(timeStart), filename ) )
}


// This functions waits for all the other goroutines to finish (when each
// one has called Done()).
// When they're all finished, we close the channel so we don't get a deadlock.
func chanMonitor( wg *sync.WaitGroup, ch chan Word ){
	wg.Wait()
	close( ch )
}

func main( ) {

	if( len(os.Args) < 2 ){
		fmt.Fprintf( os.Stderr, "Expected at least one filename.\n" )
		os.Exit(1)
	}

	filenames := os.Args[1:]

	var wordCountResult = make( chan Word )


	// **** WaitGroups ****
	// A WaitGroup is useful when you have a single channel shared amongst
	// goroutines. How do you know when all the goroutines are done with it?
	//
	// When you start a goroutine that uses a shared channel, you Add() to the WaitGroup.
	// Once processing is finished, it calls Done(). Somewhere you will have to Wait() until
	// all the Add()s have a corresponding Done(). Then the Wait() will unblock and you can
	// close your shared channel.
	wg := sync.WaitGroup{}

	// Set the start time so we can measure performance.
	timeStart := time.Now()

	for i := 0; i < len( filenames ); i++ {

		// **** WaitGroups ***
		// wordCounter uses a shared output channel. We must use a WaitGroup
		// in order to close it only when each routine is Done() with it.
		wg.Add(1)

		// Execute wordCounter. We pass it the WaitGroup to send Done() to, the filename we're processing,
		// the output channel to send the results to.
		go wordCounter( filenames[i], &wg, wordCountResult )

	}

	// The wordCounters share an output channel so instead of the routines closing the 
	// channel themselves once finished, we have to wait until they all say "I'm Done()"
	// Once all the routines are Done(), the routine closes and unblocks the channel.
	go chanMonitor( &wg, wordCountResult ) // This blocks 

	// Receive data from the results channel and output it.
	for result := range wordCountResult {
		logLine( os.Stdout, fmt.Sprintf( "%d\twords in %s", result.words, result.fileName ) )
	}

	// Show total elapsed time.
	timeEnd := time.Now()
	logLine( os.Stderr, fmt.Sprintf( "Total elapsed time %s\n", timeEnd.Sub( timeStart ) ) )

}

