package demo4

import (
	"io"
	"os"
	"fmt"
	"unicode"
	"log"
	"time"
	"sync"
	)

// Our FileWords type.
type FileWord struct {
	FileName string
	Word string
}

// Our buffer size
const chunkSize int32 = 32*1024


// This function opens a text file, reads the data in chunks and extracts words from it.
// It sends the words (of type Word) to an output channel (out)
func WordReader( filename string, out chan <- FileWord ) {

	// Call Done() when this function ends
	defer close( out )
	
	// Start a timer
	timeStart := time.Now()

	log.Output( 2, fmt.Sprintf( "wordReader started for %s", filename ) )

	// Attempt to open the file
	file, err := os.Open( filename )
	defer file.Close()

	if( err != nil ){
		log.Output( 2, fmt.Sprintf( "Error opening file %s", filename ) )
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
						out <- FileWord{ file.Name(), lastWord }	// Send the word to the output channel
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

	file.Close()
	timeEnd := time.Now()
	log.Output( 2, fmt.Sprintf( "wordReader finished in %s for %s", timeEnd.Sub(timeStart), filename ) )
}

// This function receives data from an input channel (wordChan) containing words.
// It counts the number of words and writes the results to the output channel (resultChan)
func WordCounter( wg *sync.WaitGroup, filename string, wordChan <- chan FileWord, resultChan chan <- string ) {
	
	// Call Done() when this function exits
	defer wg.Done()

	timeStart := time.Now()
	log.Output( 2, fmt.Sprintf( "wordCounter started for %s", filename ) )
	
	// Create a map containing the filename and the number of words
	var wordCountMap map[string]int32 = make( map[string]int32  )

	// Receive words from the input channel and aggregate
	for word := range wordChan {
		wordCountMap[word.FileName] += 1
	}
	
	// Send results to the output channel
	for file, words := range wordCountMap {
		resultChan <- fmt.Sprintf( "%d\twords in %s", words, file )
	}

	timeEnd := time.Now()
	log.Output( 2,  fmt.Sprintf( "wordCounter finished in %s for %s", timeEnd.Sub( timeStart ), filename ) )
}