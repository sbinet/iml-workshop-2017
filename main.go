package main

import (
	"flag"
	"image/color"
	"log"
	"os"

	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"go-hep.org/x/hep/rootio"
)

func main() {
	log.SetPrefix("iml: ")
	log.SetFlags(0)

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	fname := flag.Arg(0)
	f, err := rootio.Open(fname)
	if err != nil {
		log.Fatal(err)
	}

	obj, ok := f.Get("treeJets")
	if !ok {
		log.Fatal("no treeJets")
	}
	tree := obj.(rootio.Tree)

	log.Printf("tree %q nevts=%d\n", tree.Name(), tree.Entries())

	sc, err := rootio.NewScanner(tree, &Event{})
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	var (
		hNtowers = hbook.NewH1D(50, 0, 50)
		hMass    = hbook.NewH1D(100, 0, 50)
	)

	for sc.Next() {
		var evt Event
		err = sc.Scan(&evt)
		if err != nil {
			log.Fatal(err)
		}
		if sc.Entry()%(tree.Entries()/10) == 0 {
			log.Printf("evt[%d] %#v\n", sc.Entry(), evt)
		}
		hNtowers.Fill(float64(evt.NTowers), 1)
		hMass.Fill(float64(evt.JetMass), 1)
	}

	err = sc.Err()
	if err != nil {
		log.Fatal(err)
	}

	{
		plot, err := hplot.NewTiledPlot(draw.Tiles{Cols: 2, Rows: 1})
		if err != nil {
			log.Fatal(err)
		}
		h1, err := hplot.NewH1D(hNtowers)
		if err != nil {
			log.Fatal(err)
		}
		h1.Infos.Style = hplot.HInfoSummary
		h1.Color = color.RGBA{255, 0, 0, 255}

		h2, err := hplot.NewH1D(hMass)
		if err != nil {
			log.Fatal(err)
		}
		h2.Infos.Style = hplot.HInfoSummary
		h2.Color = color.RGBA{0, 0, 255, 255}

		p := plot.Plot(0, 0)
		p.Title.Text = "Towers multiplicity"
		p.X.Label.Text = "# of towers"
		p.Add(h1)
		p.Add(hplot.NewGrid())

		p = plot.Plot(0, 1)
		p.Title.Text = "Reconstructed Jet Mass"
		p.X.Label.Text = "Jet Mass [GeV/c^2]"
		p.Add(h2)
		p.Add(hplot.NewGrid())

		err = plot.Save(15*vg.Centimeter, -1, "plots.png")
		if err != nil {
			log.Fatal(err)
		}
	}
}

type Event struct {
	JetPt       float32   `rootio:"jetPt"`
	JetEta      float32   `rootio:"jetEta"`
	JetPhi      float32   `rootio:"jetPhi"`
	JetMass     float32   `rootio:"jetMass"`
	NTracks     int32     `rootio:"ntracks"`
	NTowers     int32     `rootio:"ntowers"`
	TrackPt     []float32 `rootio:"trackPt"`
	TrackEta    []float32 `rootio:"trackEta"`
	TrackPhi    []float32 `rootio:"trackPhi"`
	TrackCharge []float32 `rootio:"trackCharge"`
	TowerEne    []float32 `rootio:"towerE"`
	TowerEem    []float32 `rootio:"towerEem"`
	TowerEhad   []float32 `rootio:"towerEhad"`
	TowerEta    []float32 `rootio:"towerEta"`
	TowerPhi    []float32 `rootio:"towerPhi"`
}
