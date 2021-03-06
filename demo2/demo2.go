package demo2

import (
	"io"
	"os"
	"fmt"
	"unicode"
	"log"
	"time"
	"sync"
	)

// Our buffer size
const chunkSize int32 = 32*1024

// This function opens a text file, reads the data in chunks and extracts words from it.
// It outputs the word count once done.
func WordCounter( filename string, wg *sync.WaitGroup) {

	defer wg.Done()

	// Start a timer
	timeStart := time.Now()
	
	log.Output( 2, fmt.Sprintf( "wordReader started for %s", filename ) )

	var totalBytesRead, wordsRead int64	// Bytes read
	var str, lastWord string	// str is an individual character. lastWord is a word.
	
	// Attempt to open the file
	file, err := os.Open( filename )
	defer file.Close()

	if( err != nil ){
		log.Output( 2, fmt.Sprintf( "Error opening file %s", filename ) )
		return
	}


	// Start processing data from the file
	for {

		buf := make( []byte, chunkSize )	// Create a buffer
		bytes, err := file.Read( buf )		// Read data into the buffer

		// Count words only if we read some data.

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
			log.Output( 2, fmt.Sprintf( "Error reading file %s: %s",  filename, err ) )
			break
		}

	}

	timeEnd := time.Now()
	log.Output( 2, fmt.Sprintf( "%d\twords in %s (%s)", wordsRead, filename, timeEnd.Sub(timeStart) ) )
}