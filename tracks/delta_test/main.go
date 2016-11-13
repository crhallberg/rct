package main

import (
  "bufio"
  "fmt"
  "os"

	"github.com/kevinburke/rct/genetic"
	"github.com/kevinburke/rct/geo"
	"github.com/kevinburke/rct/rle"
	"github.com/kevinburke/rct/td6"
	"github.com/kevinburke/rct/tracks"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {

  fmt.Println()
  fmt.Println("Small turns test")
  fmt.Println("Expected {-1, 0, 0}, 0")
  trackPieces := genetic.CreateStation()
  trackPieces = append(trackPieces, tracks.Element{ Segment: tracks.TS_MAP[tracks.ELEM_RIGHT_QUARTER_TURN_3_TILES] })
  trackPieces = append(trackPieces, tracks.Element{ Segment: tracks.TS_MAP[tracks.ELEM_LEFT_QUARTER_TURN_3_TILES] })
  trackPieces = append(trackPieces, tracks.Element{ Segment: tracks.TS_MAP[tracks.ELEM_LEFT_QUARTER_TURN_3_TILES] })
  trackPieces = append(trackPieces, tracks.Element{ Segment: tracks.TS_MAP[tracks.ELEM_LEFT_QUARTER_TURN_3_TILES] })
  trackPieces = append(trackPieces, tracks.Element{ Segment: tracks.TS_MAP[tracks.ELEM_RIGHT_QUARTER_TURN_3_TILES] })
  trackPieces = append(trackPieces, tracks.Element{ Segment: tracks.TS_MAP[tracks.ELEM_LEFT_QUARTER_TURN_3_TILES] })
  trackPieces = append(trackPieces, tracks.Element{ Segment: tracks.TS_MAP[tracks.ELEM_RIGHT_QUARTER_TURN_3_TILES] })
  trackPieces = append(trackPieces, tracks.Element{ Segment: tracks.TS_MAP[tracks.ELEM_LEFT_QUARTER_TURN_3_TILES] })
  trackPieces = append(trackPieces, tracks.Element{ Segment: tracks.TS_MAP[tracks.ELEM_LEFT_QUARTER_TURN_3_TILES] })
  trackPieces = append(trackPieces, tracks.Element{ Segment: tracks.TS_MAP[tracks.ELEM_RIGHT_QUARTER_TURN_3_TILES] })
  trackPieces = append(trackPieces, tracks.Element{ Segment: tracks.TS_MAP[tracks.ELEM_LEFT_QUARTER_TURN_3_TILES] })
  trackPieces = append(trackPieces, tracks.Element{ Segment: tracks.TS_MAP[tracks.ELEM_LEFT_QUARTER_TURN_3_TILES] })
  writeToFile(trackPieces, "tracks/delta_test/vector_test.td6")

  fmt.Println()
  fmt.Println("S_BEND test")
  fmt.Println("Expected {14, 3, 0}, 0")
  trackPieces = genetic.CreateStation()
  trackPieces = append(trackPieces, tracks.Element{ Segment: tracks.TS_MAP[tracks.ELEM_S_BEND_LEFT] })
  trackPieces = append(trackPieces, tracks.Element{ Segment: tracks.TS_MAP[tracks.ELEM_S_BEND_LEFT] })
  trackPieces = append(trackPieces, tracks.Element{ Segment: tracks.TS_MAP[tracks.ELEM_S_BEND_LEFT] })
  writeToFile(trackPieces, "tracks/delta_test/vector_s_test.td6")

  fmt.Println()
}

func writeToFile(trackPieces []tracks.Element, filename string) {
  // Complete ride
  simple := td6.CreateMineTrainRide(trackPieces, false)
  // Print trackEnd vector
  vectors := geo.Vectors(trackPieces)
  trackEnd := vectors[len(trackPieces)-1]
  fmt.Printf("%#v\n", trackEnd)
  /*
  // The last vector does not include the changes of the last piece
  // Applying the last track piece gets you closer: fixes direction, etc
  trackEnd = geo.AdvanceVector(trackEnd, trackPieces[len(trackPieces)-1].Segment
  */
  // Convert to format
  bits, err := td6.Marshal(simple)
  check(err)
  paddedBits := td6.Pad(bits)
  // Write file
  f, err := os.Create(filename)
  check(err);
  w := bufio.NewWriter(f)
  rleWriter := rle.NewWriter(w)
  rleWriter.Write(paddedBits)
  w.Flush()
}
