/*

This program counts words in text files, given as command line parameters.

There is concurrency using goroutines, and a channel for the results of word counting.

Author: Jamie Skipworth

*/
package main

import (
	"fmt"
	"os"
	"time"
	"sync"
	"./demo3"
	"log"
	)


// This functions waits for all the other goroutines to finish (when each
// one has called Done()).
// When they're all finished, we close the channel so we don't get a deadlock.
func chanMonitor( wg *sync.WaitGroup, ch chan demo3.FileWords ){
	wg.Wait()
	close( ch )
}

func main( ) {

	// Set the format of the logger
	log.SetFlags( log.Ldate|log.Lmicroseconds )

	if( len(os.Args) < 2 ){
		log.Output( 2, "Expected at least one filename.\n" )
		os.Exit(1)
	}

	filenames := os.Args[1:]

	var wordCountResult = make( chan demo3.FileWords )

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
		go demo3.WordCounter( filenames[i], &wg, wordCountResult )

	}

	// The wordCounters share an output channel so instead of the routines closing the 
	// channel themselves once finished, we have to wait until they all say "I'm Done()"
	// Once all the routines are Done(), the routine closes and unblocks the channel.
	go chanMonitor( &wg, wordCountResult ) // This blocks 

	// Receive data from the results channel and output it.
	for result := range wordCountResult {
		log.Output( 2, fmt.Sprintf( "%d\twords in %s", result.Words, result.FileName ) )
	}

	// Show total elapsed time.
	timeEnd := time.Now()
	log.Output( 2, fmt.Sprintf( "Total elapsed time %s\n", timeEnd.Sub( timeStart ) ) )

}

