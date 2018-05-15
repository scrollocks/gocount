/*

This program counts words in text files, given as command line parameters.

There is a wordReader for extracting words from text files, which are sent to a words channel.
There is a wordCounter that receives words from the words channel, and sends results to the results channel.

Author: Jamie Skipworth

*/
package main

import (
	"fmt"
	"os"
	"time"
	"sync"
	"./demo4"
	"log"
	)


// This functions waits for goroutines to finish.
// When they're all finished, we close the given channel so we don't get a deadlock.
func chanMonitor( wg *sync.WaitGroup, ch chan string ){
	wg.Wait()
	close( ch )
}

// Count some words!
func main( ) {

	// Set the format of the logger
	log.SetFlags( log.Ldate|log.Lmicroseconds )
	
	// Make sure we have at least one cmd line parameter
	if( len(os.Args) < 2 ){
		log.Output( 2,  fmt.Sprintf( "Expected at least one filename." ) )
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
		wordChan := make( chan demo4.FileWord )

		// Execute wordReader on the filename, and give it a channel to send words to
		go demo4.WordReader( filenames[i], wordChan )

		// **** WaitGroups ***
		// wordCounter uses a shared output channel. We must use a WaitGroup
		// in order to close it only when each routine is Done() with it.
		wg.Add( 1 ) 

		// Execute wordCounter. We pass it the WaitGroup to send Done() to, the filename we're processing,
		// the input channel of words, and the channel to send the results to.
		go demo4.WordCounter( &wg, filenames[i], wordChan, resultChan )

	}

	// The wordCounters share an output channel so instead of the routines closing the 
	// channel themselves once finished, we have to wait until they all say "I'm Done()"
	// Once all the routines are Done(), the routine closes and unblocks the channel.
	go chanMonitor( &wg, resultChan )

	// Receive data from the results channel and output it.
	for countResults := range resultChan {
		log.Output( 2, countResults )
	}

	// Show total elapsed time.
	timeEnd := time.Now()
	log.Output( 2, fmt.Sprintf( "Total elapsed time %s", timeEnd.Sub( timeStart )) )


}

