package demo3

import (
	"io"
	"os"
	"fmt"
	"unicode"
	"log"
	"time"
	"sync"
	)

const chunkSize int32 = 32*1024

// Our FileWords type.
type FileWords struct {
	FileName string
	Words int64
}

// This function opens a text file, reads the data in chunks and extracts words from it.
// It sends the word to a results channel.
func WordCounter( filename string, waitGroup *sync.WaitGroup, out chan <- FileWords ) {

	defer waitGroup.Done()

	// Call Done() when this function ends
	timeStart := time.Now()

	log.Output( 2, fmt.Sprintf( "wordReader started for %s", filename ) )

	// Attempt to open the file
	file, err := os.Open( filename )
	defer file.Close()

	if( err != nil ){
		log.Output( 2, fmt.Sprintf( "Error opening file %s", filename ) )
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
			log.Output( 2, fmt.Sprintf( "Error reading file %s: %s",  filename, err ) )
			break
		}

	}

	out <- FileWords{ file.Name(), wordCount }
	file.Close()
	timeEnd := time.Now()
	log.Output( 2, fmt.Sprintf( "wordCounter finished in %s for %s", timeEnd.Sub(timeStart), filename ) )
}