/*

This program counts words in text files, given as command line parameters.

There is a wordReader for extracting words from text files, which are sent to a words channel.
There is a wordCounter that receives words from the words channel, and sends results to the results channel.

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

// A struct to encapsulate a filename and a word
type Word struct {
	fileName string
	word string
}

// Our buffer size
const chunkSize int32 = 32*1024

// Writes text to a file descriptor. Nice for consistent logging.
func logLine( fd *os.File, message string ) {
	fmt.Fprintf( fd, "%s - %s\n", time.Now().Format("20060106 15:04:05.000"), message )
}

// This function opens a text file, reads the data in chunks and extracts words from it.
// It sends the words (of type Word) to an output channel (out)
func wordReader( filename string, out chan <- Word ) {

	// Call Done() when this function ends
	defer close( out )
	
	// Start a timer
	timeStart := time.Now()

	logLine( os.Stderr, fmt.Sprintf( "wordReader started for %s", filename ) )

	// Attempt to open the file
	file, err := os.Open( filename )
	defer file.Close()

	if( err != nil ){
		logLine( os.Stderr, fmt.Sprintf( "Error opening file %s", filename ) )
		return
	}

	var totalBytesRead int64	// Bytes read
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
						out <- Word{ file.Name(), lastWord }	// Send the word to the output channel
						lastWord = ""	// Blank out the word ready for the next one
					}
				// If the run is not whitespace then add it to the lastWord string. We're not in a region of white-space.
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

	file.Close()
	timeEnd := time.Now()
	logLine( os.Stderr, fmt.Sprintf( "wordReader finished in %s for %s", timeEnd.Sub(timeStart), filename ) )
}

// This function receives data from an input channel (wordChan) containing words.
// It counts the number of words and writes the results to the output channel (resultChan)
func wordCounter( wg *sync.WaitGroup, filename string, wordChan <- chan Word, resultChan chan <- string ) {
	
	// Call Done() when this function exits
	defer wg.Done()

	timeStart := time.Now()
	logLine( os.Stderr, fmt.Sprintf( "wordCounter started for %s", filename ) )
	
	// Create a map containing the filename and the number of words
	var wordCountMap map[string]int32 = make( map[string]int32  )

	// Receive words from the input channel and aggregate
	for word := range wordChan {
		wordCountMap[word.fileName] += 1
	}
	
	// Send results to the output channel
	for file, words := range wordCountMap {
		resultChan <- fmt.Sprintf( "%d\twords in %s", words, file )
	}

	timeEnd := time.Now()
	logLine( os.Stderr, fmt.Sprintf( "wordCounter finished in %s for %s", timeEnd.Sub( timeStart ), filename ) )
}

// This functions waits for goroutines to finish.
// When they're all finished, we close the given channel so we don't get a deadlock.
func chanMonitor( wg *sync.WaitGroup, ch chan string ){
	wg.Wait()
	close( ch )
}

// Count some words!
func main( ) {

	// Make sure we have at least one cmd line parameter
	if( len(os.Args) < 2 ){
		logLine( os.Stderr, fmt.Sprintf( "Expected at least one filename." ) )
		os.Exit(1)
	}

	// Take a slice of the command line. Arg[0] is the program name.
	filenames := os.Args[1:]

	// Create a channel that our results will be sent to.
	resultChan := make( chan string )

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

	// Iterate through files from the command line arguments
	for i := 0; i < len( filenames ); i++ {

		// For each file, create a channel into which we'll send words.
		wordChan := make( chan Word )

		// Execute wordReader on the filename, and give it a channel to send words to
		go wordReader( filenames[i], wordChan )

		// **** WaitGroups ***
		// wordCounter uses a shared output channel. We must use a WaitGroup
		// in order to close it only when each routine is Done() with it.
		wg.Add( 1 ) 

		// Execute wordCounter. We pass it the WaitGroup to send Done() to, the filename we're processing,
		// the input channel of words, and the channel to send the results to.
		go wordCounter( &wg, filenames[i], wordChan, resultChan )

	}

	// The wordCounters share an output channel so instead of the routines closing the 
	// channel themselves once finished, we have to wait until they all say "I'm Done()"
	// Once all the routines are Done(), the routine closes and unblocks the channel.
	go chanMonitor( &wg, resultChan )

	// Receive data from the results channel and output it.
	for countResults := range resultChan {
		logLine( os.Stdout, countResults )
	}

	// Show total elapsed time.
	timeEnd := time.Now()
	logLine( os.Stdout, fmt.Sprintf( "Total elapsed time %s", timeEnd.Sub( timeStart )) )


}

