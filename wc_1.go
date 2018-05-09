/*

This program counts words in text files, given as command line parameters.

There is no concurrency, nor any channels

Author: Jamie Skipworth

*/
package main
import (
	"fmt"
	"io"
	"os"
	"unicode"
	"time"
	)


const chunkSize int32 = 32*1024

// Writes text to a file descriptor. Nice for consistent logging.
func logLine( fd *os.File, message string ) {
	fmt.Fprintf( fd, "%s - %s\n", time.Now().Format("20060106 15:04:05.000"), message )
}

// This function opens a text file, reads the data in chunks and extracts words from it.
// It outputs the word count once done.
func wordCounter( filename string ) ( bytes int64, words int64 ) {

	timeStart := time.Now()

	var bytesRead, wordsRead int64 
	var str, lastWord string

	logLine( os.Stderr, fmt.Sprintf( "wordCounter started for %s", filename ) )
	
	// Attempt to open the file
	file, err := os.Open( filename )
	defer file.Close()

	if( err != nil ){
		logLine( os.Stderr, fmt.Sprintf( "Error opening file %s", filename ) )
		return
	}

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
						wordsRead += 1	// Increment the word counter
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

	timeEnd := time.Now()
	logLine( os.Stderr, fmt.Sprintf( "wordCounter finished in %s for %s", timeEnd.Sub(timeStart), filename ) )
	return bytesRead, wordsRead

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

	// Set the start time so we can measure performance.
	timeStart := time.Now()

	// Iterate through files from the command line arguments
	for i := 0; i < len( filenames ); i++ {

		// Call the wordCounter. Discard bytesRead (1st return param), but keep the word count (2nd)
		_, w := wordCounter( filenames[i] )
		logLine( os.Stdout, fmt.Sprintf( "%d\twords in %s", w, filenames[i] ) )
	}

	// Show total elapsed time.
	timeEnd := time.Now()
	logLine( os.Stdout, fmt.Sprintf( "Total elapsed time %s", timeEnd.Sub( timeStart )) )


}
