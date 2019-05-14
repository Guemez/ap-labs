package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type Node struct {
	id        int
	neighbors []*Node
	history   []*Node
	vec       pixel.Vec
	isEaten   bool
}

type wall struct {
	rect pixel.Rect
}

const constSpeed = 150.0

var imd = imdraw.New(nil)
var speed = constSpeed
var pause = false
var score = 0
var numGhosts int
var currentPos pixel.Vec
var nodes = make(map[int]*Node, 250)
var currentNodePacman *Node
var walls = []wall{}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open("./assets/" + path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func checkCollision(vec1, vec2 pixel.Vec) bool { //Check if vec1 is in boundary of vec2
	if vec1.X > vec2.X-20 && vec1.X < vec2.X+20 {
		if vec1.Y > vec2.Y-20 && vec1.Y < vec2.Y+20 {
			return true
		}
	}
	return false
}

func checkCollisionWall(vec pixel.Vec, rect pixel.Rect) bool { //Check if vec1 is colliding with the rectangle using bounding box algorithm
	newRect := pixel.R(vec.X-15, vec.Y-15, vec.X+15, vec.Y+15)
	if newRect.Min.X < rect.Min.X+rect.W() && newRect.Min.X+newRect.W() > rect.Min.X && newRect.Min.Y < rect.Min.Y+rect.H() && newRect.Min.Y+newRect.H() > rect.Min.Y {
		// collision detected!
		return true
	}
	return false
}

func checkCollisionAllWalls(vec pixel.Vec) bool {
	for _, singeWall := range walls {
		if checkCollisionWall(vec, singeWall.rect) {
			return true
		}
	}
	return false
}

func checkCollisionGhost(vec pixel.Vec) bool {
	if checkCollision(vec, currentPos) {
		log.Printf("YOU LOSE!!")
		os.Exit(0)
	}
	return false
}

func printNode(n Node) {
	fmt.Printf("Current ID: %d\n", n.id)
	fmt.Print("Neighbors: \n")
	for i := range n.neighbors {
		fmt.Print("ID: ")
		fmt.Println(n.neighbors[i].id)
	}

}

func setUpLevel() {
	{
		imd.Color = colornames.Blue
		limits := []wall{
			{rect: pixel.R(0, 0, 1024, 720)}, //Border
		}
		for _, l := range limits {
			imd.Push(l.rect.Min, l.rect.Max)
		}
		imd.Rectangle(20)
		imd.Color = colornames.Black
		exits := []wall{
			{rect: pixel.R(-10, 360, 0, 360+30)},        //Left Exit
			{rect: pixel.R(1000, 360-30, 1024, 360+30)}, //Right Exit
		}
		for _, e := range exits {
			imd.Push(e.rect.Min, e.rect.Max)
		}
		imd.Rectangle(0)
		imd.Color = colornames.Blue
		walls = []wall{
			{rect: pixel.R(70, 70, 300, 300)},                   //SouthWest block
			{rect: pixel.R(1024-300, 70, 1024-70, 300)},         //SouthEast block
			{rect: pixel.R(70, 720-300, 300, 720-70)},           //NorthWest block
			{rect: pixel.R(1024-300, 720-300, 1024-70, 720-70)}, //NorthEast block

			{rect: pixel.R(512-150, 360-100, 512-100, 360+100)},
			{rect: pixel.R(512-150, 360-100, 512+150, 360-50)},
			{rect: pixel.R(512+100, 360-100, 512+150, 360+100)}, //Center block

			{rect: pixel.R(512-50, 360+200-50, 512+50, 360+300+50)}, //Center North block
			{rect: pixel.R(512-50, 360-300-50, 512+50, 360-200+50)}, //Center South block
			{rect: pixel.R(512+325-25, 360-25, 512+325+25, 360+25)}, //Center East block
			{rect: pixel.R(512-325-25, 360-25, 512-325+25, 360+25)}, //Center West block
		}
		for _, w := range walls {
			imd.Push(w.rect.Min, w.rect.Max)
			imd.Rectangle(0)
		}

		//Declare nodes
		{
			{ //row 1
				nodes[0] = &Node{0, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(40)), false}
				nodes[1] = &Node{1, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(80), float64(40)), false}
				nodes[2] = &Node{2, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(120), float64(40)), false}
				nodes[3] = &Node{3, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(160), float64(40)), false}
				nodes[4] = &Node{4, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(200), float64(40)), false}
				nodes[5] = &Node{5, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(240), float64(40)), false}
				nodes[6] = &Node{6, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(280), float64(40)), false}
				nodes[7] = &Node{7, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(40)), false}
				nodes[8] = &Node{8, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(360), float64(40)), false}
				nodes[9] = &Node{9, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(400), float64(40)), false}
				nodes[10] = &Node{10, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(40)), false}
				nodes[11] = &Node{11, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(40)), false}
				nodes[12] = &Node{12, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(620), float64(40)), false}
				nodes[13] = &Node{13, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(660), float64(40)), false}
				nodes[14] = &Node{14, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(40)), false}
				nodes[15] = &Node{15, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(740), float64(40)), false}
				nodes[16] = &Node{16, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(780), float64(40)), false}
				nodes[17] = &Node{17, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(820), float64(40)), false}
				nodes[18] = &Node{18, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(860), float64(40)), false}
				nodes[19] = &Node{19, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(900), float64(40)), false}
				nodes[20] = &Node{20, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(940), float64(40)), false}
				nodes[21] = &Node{21, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(40)), false}

				/*for i := 0; i < 11; i++ {
					nodes[i] = &Node{0, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(0+40*(i+1)), float64(40)), false}
				, false}
				for i := 12; i < 22; i++ {
					nodes[i] = &Node{0, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580+(40*(i-12))), float64(40)), false}
				, false}*/
			}
			{ //row2
				nodes[22] = &Node{22, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(80)), false}
				nodes[23] = &Node{23, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(80)), false}
				nodes[24] = &Node{24, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(360), float64(80)), false}
				nodes[25] = &Node{25, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(400), float64(80)), false}
				nodes[26] = &Node{26, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(80)), false}
				nodes[27] = &Node{27, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(80)), false}
				nodes[28] = &Node{28, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(620), float64(80)), false}
				nodes[29] = &Node{29, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(660), float64(80)), false}
				nodes[30] = &Node{30, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(80)), false}
				nodes[31] = &Node{31, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(80)), false}
			}
			{ //row3
				nodes[32] = &Node{32, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(120)), false}
				nodes[33] = &Node{33, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(120)), false}
				nodes[34] = &Node{34, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(360), float64(120)), false}
				nodes[35] = &Node{35, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(400), float64(120)), false}
				nodes[36] = &Node{36, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(120)), false}
				nodes[37] = &Node{37, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(120)), false}
				nodes[38] = &Node{38, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(620), float64(120)), false}
				nodes[39] = &Node{39, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(660), float64(120)), false}
				nodes[40] = &Node{40, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(120)), false}
				nodes[41] = &Node{41, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(120)), false}
			}
			{ //row4
				nodes[42] = &Node{42, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(160)), false}
				nodes[43] = &Node{43, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(160)), false}
				nodes[44] = &Node{44, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(360), float64(160)), false}
				nodes[45] = &Node{45, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(400), float64(160)), false}
				nodes[46] = &Node{46, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(160)), false}
				nodes[47] = &Node{47, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(160)), false}
				nodes[48] = &Node{48, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(620), float64(160)), false}
				nodes[49] = &Node{49, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(660), float64(160)), false}
				nodes[50] = &Node{50, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(160)), false}
				nodes[51] = &Node{51, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(160)), false}
			}
			{ //row5
				nodes[52] = &Node{52, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(200)), false}
				nodes[53] = &Node{53, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(200)), false}
				nodes[54] = &Node{54, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(360), float64(200)), false}
				nodes[55] = &Node{55, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(400), float64(200)), false}
				nodes[56] = &Node{56, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(200)), false}
				nodes[57] = &Node{57, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(200)), false}
				nodes[58] = &Node{58, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(620), float64(200)), false}
				nodes[59] = &Node{59, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(660), float64(200)), false}
				nodes[60] = &Node{60, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(200)), false}
				nodes[61] = &Node{61, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(200)), false}
			}
			{ //row6
				nodes[62] = &Node{62, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(240)), false}
				nodes[63] = &Node{63, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(240)), false}
				nodes[64] = &Node{64, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(360), float64(240)), false}
				nodes[65] = &Node{65, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(400), float64(240)), false}
				nodes[66] = &Node{66, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(240)), false}
				nodes[67] = &Node{72, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(485), float64(240)), false}
				nodes[68] = &Node{73, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(535), float64(240)), false}
				nodes[69] = &Node{67, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(240)), false}
				nodes[70] = &Node{68, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(620), float64(240)), false}
				nodes[71] = &Node{69, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(660), float64(240)), false}
				nodes[72] = &Node{70, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(240)), false}
				nodes[73] = &Node{71, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(240)), false}
			}
			{ //row7
				nodes[74] = &Node{74, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(280)), false}
				nodes[75] = &Node{75, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(280)), false}
				nodes[76] = &Node{76, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(280)), false}
				nodes[77] = &Node{77, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(280)), false}
			}
			{ //row8
				nodes[78] = &Node{78, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(320)), false}
				nodes[79] = &Node{79, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(80), float64(320)), false}
				nodes[80] = &Node{80, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(120), float64(320)), false}
				nodes[81] = &Node{81, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(160), float64(320)), false}
				nodes[82] = &Node{82, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(200), float64(320)), false}
				nodes[83] = &Node{83, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(240), float64(320)), false}
				nodes[84] = &Node{84, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(280), float64(320)), false}
				nodes[85] = &Node{85, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(320)), false}
				nodes[86] = &Node{86, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(320)), false}
				nodes[87] = &Node{87, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(740), float64(320)), false}
				nodes[88] = &Node{88, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(780), float64(320)), false}
				nodes[89] = &Node{89, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(820), float64(320)), false}
				nodes[90] = &Node{90, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(860), float64(320)), false}
				nodes[91] = &Node{91, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(900), float64(320)), false}
				nodes[92] = &Node{92, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(940), float64(320)), false}
				nodes[93] = &Node{93, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(320)), false}
			}
			{ //row9
				nodes[94] = &Node{94, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(360)), false}
				nodes[95] = &Node{95, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(80), float64(360)), false}
				nodes[96] = &Node{96, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(120), float64(360)), false}
				nodes[97] = &Node{97, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(240), float64(360)), false}
				nodes[98] = &Node{98, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(280), float64(360)), false}
				nodes[99] = &Node{99, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(360)), false}
				nodes[100] = &Node{100, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(360)), false}
				nodes[101] = &Node{101, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(740), float64(360)), false}
				nodes[102] = &Node{102, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(780), float64(360)), false}
				nodes[103] = &Node{103, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(900), float64(360)), false}
				nodes[104] = &Node{104, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(940), float64(360)), false}
				nodes[105] = &Node{105, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(360)), false}
			}
			{ //row10
				nodes[106] = &Node{106, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(400)), false}
				nodes[107] = &Node{107, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(80), float64(400)), false}
				nodes[108] = &Node{108, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(120), float64(400)), false}
				nodes[109] = &Node{109, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(160), float64(400)), false}
				nodes[110] = &Node{110, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(200), float64(400)), false}
				nodes[111] = &Node{111, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(240), float64(400)), false}
				nodes[112] = &Node{112, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(280), float64(400)), false}
				nodes[113] = &Node{113, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(400)), false}
				nodes[114] = &Node{114, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(400)), false}
				nodes[115] = &Node{115, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(740), float64(400)), false}
				nodes[116] = &Node{116, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(780), float64(400)), false}
				nodes[117] = &Node{117, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(820), float64(400)), false}
				nodes[118] = &Node{118, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(860), float64(400)), false}
				nodes[119] = &Node{119, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(900), float64(400)), false}
				nodes[120] = &Node{120, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(940), float64(400)), false}
				nodes[121] = &Node{121, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(400)), false}
			}
			{ //row11
				nodes[122] = &Node{122, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(440)), false}
				nodes[123] = &Node{123, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(440)), false}
				nodes[124] = &Node{124, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(440)), false}
				nodes[125] = &Node{125, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(440)), false}
			}
			{ //row12
				nodes[126] = &Node{126, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(480)), false}
				nodes[127] = &Node{127, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(480)), false}
				nodes[128] = &Node{128, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(360), float64(480)), false}
				nodes[129] = &Node{129, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(400), float64(480)), false}
				nodes[130] = &Node{130, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(480)), false}
				nodes[131] = &Node{131, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(485), float64(480)), false}
				nodes[132] = &Node{132, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(535), float64(480)), false}
				nodes[133] = &Node{133, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(480)), false}
				nodes[134] = &Node{134, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(620), float64(480)), false}
				nodes[135] = &Node{135, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(660), float64(480)), false}
				nodes[136] = &Node{136, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(480)), false}
				nodes[137] = &Node{137, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(480)), false}
			}
			{ //row13
				nodes[138] = &Node{138, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(520)), false}
				nodes[139] = &Node{139, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(520)), false}
				nodes[140] = &Node{140, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(360), float64(520)), false}
				nodes[141] = &Node{141, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(400), float64(520)), false}
				nodes[142] = &Node{142, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(520)), false}
				nodes[143] = &Node{143, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(520)), false}
				nodes[144] = &Node{144, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(620), float64(520)), false}
				nodes[145] = &Node{145, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(660), float64(520)), false}
				nodes[146] = &Node{146, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(520)), false}
				nodes[147] = &Node{147, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(520)), false}
			}
			{ //row14
				nodes[148] = &Node{148, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(560)), false}
				nodes[149] = &Node{149, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(560)), false}
				nodes[150] = &Node{150, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(360), float64(560)), false}
				nodes[151] = &Node{151, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(400), float64(560)), false}
				nodes[152] = &Node{152, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(560)), false}
				nodes[153] = &Node{153, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(560)), false}
				nodes[154] = &Node{154, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(620), float64(560)), false}
				nodes[155] = &Node{155, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(660), float64(560)), false}
				nodes[156] = &Node{156, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(560)), false}
				nodes[157] = &Node{157, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(560)), false}
			}
			{ //row15
				nodes[158] = &Node{158, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(600)), false}
				nodes[159] = &Node{159, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(600)), false}
				nodes[160] = &Node{160, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(360), float64(600)), false}
				nodes[161] = &Node{161, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(400), float64(600)), false}
				nodes[162] = &Node{162, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(600)), false}
				nodes[163] = &Node{163, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(600)), false}
				nodes[164] = &Node{164, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(620), float64(600)), false}
				nodes[165] = &Node{165, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(660), float64(600)), false}
				nodes[166] = &Node{166, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(600)), false}
				nodes[167] = &Node{167, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(600)), false}
			}
			{ //row16
				nodes[168] = &Node{168, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(640)), false}
				nodes[169] = &Node{169, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(640)), false}
				nodes[170] = &Node{170, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(360), float64(640)), false}
				nodes[171] = &Node{171, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(400), float64(640)), false}
				nodes[172] = &Node{172, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(640)), false}
				nodes[173] = &Node{173, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(640)), false}
				nodes[174] = &Node{174, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(620), float64(640)), false}
				nodes[175] = &Node{175, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(660), float64(640)), false}
				nodes[176] = &Node{176, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(640)), false}
				nodes[177] = &Node{177, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(640)), false}
			}
			{ //row 17
				nodes[178] = &Node{178, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(680)), false}
				nodes[179] = &Node{179, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(80), float64(680)), false}
				nodes[180] = &Node{180, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(120), float64(680)), false}
				nodes[181] = &Node{181, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(160), float64(680)), false}
				nodes[182] = &Node{182, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(200), float64(680)), false}
				nodes[183] = &Node{183, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(240), float64(680)), false}
				nodes[184] = &Node{184, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(280), float64(680)), false}
				nodes[185] = &Node{185, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(680)), false}
				nodes[186] = &Node{186, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(360), float64(680)), false}
				nodes[187] = &Node{187, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(400), float64(680)), false}
				nodes[188] = &Node{188, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(680)), false}
				nodes[189] = &Node{189, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(680)), false}
				nodes[190] = &Node{190, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(620), float64(680)), false}
				nodes[191] = &Node{191, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(660), float64(680)), false}
				nodes[192] = &Node{192, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(680)), false}
				nodes[193] = &Node{193, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(740), float64(680)), false}
				nodes[194] = &Node{194, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(780), float64(680)), false}
				nodes[195] = &Node{195, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(820), float64(680)), false}
				nodes[196] = &Node{196, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(860), float64(680)), false}
				nodes[197] = &Node{197, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(900), float64(680)), false}
				nodes[198] = &Node{198, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(940), float64(680)), false}
				nodes[199] = &Node{199, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(680)), false}
			}
			{ //Inside Square
				nodes[200] = &Node{200, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(320)), true}
				nodes[201] = &Node{201, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(485), float64(320)), true}
				nodes[202] = &Node{202, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(535), float64(320)), true}
				nodes[203] = &Node{203, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(320)), true}

				nodes[204] = &Node{204, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(360)), true}
				nodes[205] = &Node{205, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(485), float64(360)), true}
				nodes[206] = &Node{206, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(535), float64(360)), true}
				nodes[207] = &Node{207, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(360)), true}

				nodes[208] = &Node{208, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(400)), true}
				nodes[209] = &Node{209, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(485), float64(400)), true}
				nodes[210] = &Node{210, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(535), float64(400)), true}
				nodes[211] = &Node{211, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(400)), true}

				nodes[212] = &Node{212, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(440)), true}
				nodes[213] = &Node{213, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(485), float64(440)), true}
				nodes[214] = &Node{214, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(535), float64(440)), true}
				nodes[215] = &Node{215, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(440)), true}
			}
		}

		//Declare neighbors
		{
			nodes[0].neighbors = append(nodes[0].neighbors, nodes[1], nodes[22])
			nodes[1].neighbors = append(nodes[1].neighbors, nodes[0], nodes[2])
			nodes[2].neighbors = append(nodes[2].neighbors, nodes[1], nodes[3])
			nodes[3].neighbors = append(nodes[3].neighbors, nodes[2], nodes[4])
			nodes[4].neighbors = append(nodes[4].neighbors, nodes[3], nodes[5])
			nodes[5].neighbors = append(nodes[5].neighbors, nodes[4], nodes[6])
			nodes[6].neighbors = append(nodes[6].neighbors, nodes[5], nodes[7])
			nodes[7].neighbors = append(nodes[7].neighbors, nodes[6], nodes[8], nodes[23])
			nodes[8].neighbors = append(nodes[8].neighbors, nodes[7], nodes[9], nodes[24])
			nodes[9].neighbors = append(nodes[9].neighbors, nodes[8], nodes[10], nodes[25])
			nodes[10].neighbors = append(nodes[10].neighbors, nodes[9], nodes[26])
			nodes[11].neighbors = append(nodes[11].neighbors, nodes[12], nodes[27])
			nodes[12].neighbors = append(nodes[12].neighbors, nodes[11], nodes[13], nodes[28])
			nodes[13].neighbors = append(nodes[13].neighbors, nodes[12], nodes[14], nodes[29])
			nodes[14].neighbors = append(nodes[14].neighbors, nodes[13], nodes[15], nodes[30])
			nodes[15].neighbors = append(nodes[15].neighbors, nodes[14], nodes[16])
			nodes[16].neighbors = append(nodes[16].neighbors, nodes[15], nodes[17])
			nodes[17].neighbors = append(nodes[17].neighbors, nodes[16], nodes[18])
			nodes[18].neighbors = append(nodes[18].neighbors, nodes[17], nodes[19])
			nodes[19].neighbors = append(nodes[19].neighbors, nodes[18], nodes[20])
			nodes[20].neighbors = append(nodes[20].neighbors, nodes[19], nodes[21])
			nodes[21].neighbors = append(nodes[21].neighbors, nodes[20], nodes[31])

			nodes[22].neighbors = append(nodes[22].neighbors, nodes[0], nodes[32])
			nodes[23].neighbors = append(nodes[23].neighbors, nodes[7], nodes[24], nodes[33])
			nodes[24].neighbors = append(nodes[24].neighbors, nodes[8], nodes[23], nodes[25], nodes[34])
			nodes[25].neighbors = append(nodes[25].neighbors, nodes[9], nodes[24], nodes[26], nodes[35])
			nodes[26].neighbors = append(nodes[26].neighbors, nodes[10], nodes[25], nodes[36])
			nodes[27].neighbors = append(nodes[27].neighbors, nodes[11], nodes[28], nodes[37])
			nodes[28].neighbors = append(nodes[28].neighbors, nodes[12], nodes[27], nodes[29], nodes[38])
			nodes[29].neighbors = append(nodes[29].neighbors, nodes[13], nodes[28], nodes[30], nodes[39])
			nodes[30].neighbors = append(nodes[30].neighbors, nodes[14], nodes[29], nodes[40])
			nodes[31].neighbors = append(nodes[31].neighbors, nodes[21], nodes[41])

			nodes[32].neighbors = append(nodes[32].neighbors, nodes[22], nodes[42])
			nodes[33].neighbors = append(nodes[33].neighbors, nodes[23], nodes[34], nodes[43])
			nodes[34].neighbors = append(nodes[34].neighbors, nodes[24], nodes[33], nodes[35], nodes[44])
			nodes[35].neighbors = append(nodes[35].neighbors, nodes[25], nodes[34], nodes[36], nodes[45])
			nodes[36].neighbors = append(nodes[36].neighbors, nodes[26], nodes[35], nodes[46])
			nodes[37].neighbors = append(nodes[37].neighbors, nodes[27], nodes[38], nodes[47])
			nodes[38].neighbors = append(nodes[38].neighbors, nodes[28], nodes[37], nodes[39], nodes[48])
			nodes[39].neighbors = append(nodes[39].neighbors, nodes[29], nodes[38], nodes[40], nodes[49])
			nodes[40].neighbors = append(nodes[40].neighbors, nodes[30], nodes[39], nodes[50])
			nodes[41].neighbors = append(nodes[41].neighbors, nodes[31], nodes[51])

			nodes[42].neighbors = append(nodes[42].neighbors, nodes[32], nodes[52])
			nodes[43].neighbors = append(nodes[43].neighbors, nodes[33], nodes[44], nodes[53])
			nodes[44].neighbors = append(nodes[44].neighbors, nodes[34], nodes[43], nodes[45], nodes[54])
			nodes[45].neighbors = append(nodes[45].neighbors, nodes[35], nodes[44], nodes[46], nodes[55])
			nodes[46].neighbors = append(nodes[46].neighbors, nodes[36], nodes[45], nodes[56])
			nodes[47].neighbors = append(nodes[47].neighbors, nodes[37], nodes[48], nodes[57])
			nodes[48].neighbors = append(nodes[48].neighbors, nodes[38], nodes[47], nodes[49], nodes[58])
			nodes[49].neighbors = append(nodes[49].neighbors, nodes[39], nodes[48], nodes[50], nodes[59])
			nodes[50].neighbors = append(nodes[50].neighbors, nodes[40], nodes[49], nodes[60])
			nodes[51].neighbors = append(nodes[51].neighbors, nodes[41], nodes[61])

			nodes[52].neighbors = append(nodes[52].neighbors, nodes[42], nodes[62])
			nodes[53].neighbors = append(nodes[53].neighbors, nodes[43], nodes[54], nodes[63])
			nodes[54].neighbors = append(nodes[54].neighbors, nodes[44], nodes[53], nodes[55], nodes[64])
			nodes[55].neighbors = append(nodes[55].neighbors, nodes[45], nodes[54], nodes[56], nodes[65])
			nodes[56].neighbors = append(nodes[56].neighbors, nodes[46], nodes[55], nodes[66])
			nodes[57].neighbors = append(nodes[57].neighbors, nodes[47], nodes[58], nodes[69])
			nodes[58].neighbors = append(nodes[58].neighbors, nodes[48], nodes[57], nodes[59], nodes[70])
			nodes[59].neighbors = append(nodes[59].neighbors, nodes[49], nodes[58], nodes[60], nodes[71])
			nodes[60].neighbors = append(nodes[60].neighbors, nodes[50], nodes[59], nodes[72])
			nodes[61].neighbors = append(nodes[61].neighbors, nodes[51], nodes[73])

			nodes[62].neighbors = append(nodes[62].neighbors, nodes[52], nodes[74])
			nodes[63].neighbors = append(nodes[63].neighbors, nodes[53], nodes[64], nodes[75])
			nodes[64].neighbors = append(nodes[64].neighbors, nodes[54], nodes[63], nodes[65])
			nodes[65].neighbors = append(nodes[65].neighbors, nodes[55], nodes[64], nodes[66])
			nodes[66].neighbors = append(nodes[66].neighbors, nodes[56], nodes[65], nodes[67])
			nodes[67].neighbors = append(nodes[67].neighbors, nodes[66], nodes[68])
			nodes[68].neighbors = append(nodes[68].neighbors, nodes[67], nodes[69])
			nodes[69].neighbors = append(nodes[69].neighbors, nodes[57], nodes[68], nodes[70])
			nodes[70].neighbors = append(nodes[70].neighbors, nodes[58], nodes[69], nodes[71])
			nodes[71].neighbors = append(nodes[71].neighbors, nodes[59], nodes[70], nodes[72])
			nodes[72].neighbors = append(nodes[72].neighbors, nodes[60], nodes[71], nodes[76])
			nodes[73].neighbors = append(nodes[73].neighbors, nodes[61], nodes[77])

			nodes[74].neighbors = append(nodes[74].neighbors, nodes[62], nodes[78])
			nodes[75].neighbors = append(nodes[75].neighbors, nodes[63], nodes[85])
			nodes[76].neighbors = append(nodes[76].neighbors, nodes[72], nodes[86])
			nodes[77].neighbors = append(nodes[77].neighbors, nodes[73], nodes[93])

			nodes[78].neighbors = append(nodes[78].neighbors, nodes[74], nodes[79], nodes[94])
			nodes[79].neighbors = append(nodes[79].neighbors, nodes[78], nodes[80], nodes[95])
			nodes[80].neighbors = append(nodes[80].neighbors, nodes[79], nodes[81], nodes[96])
			nodes[81].neighbors = append(nodes[81].neighbors, nodes[80], nodes[82])
			nodes[82].neighbors = append(nodes[82].neighbors, nodes[81], nodes[83])
			nodes[83].neighbors = append(nodes[83].neighbors, nodes[82], nodes[84], nodes[97])
			nodes[84].neighbors = append(nodes[84].neighbors, nodes[83], nodes[85], nodes[98])
			nodes[85].neighbors = append(nodes[85].neighbors, nodes[75], nodes[84], nodes[99])
			nodes[86].neighbors = append(nodes[86].neighbors, nodes[76], nodes[87], nodes[100])
			nodes[87].neighbors = append(nodes[87].neighbors, nodes[86], nodes[88], nodes[101])
			nodes[88].neighbors = append(nodes[88].neighbors, nodes[87], nodes[89], nodes[102])
			nodes[89].neighbors = append(nodes[89].neighbors, nodes[88], nodes[90])
			nodes[90].neighbors = append(nodes[90].neighbors, nodes[89], nodes[91])
			nodes[91].neighbors = append(nodes[91].neighbors, nodes[90], nodes[92], nodes[103])
			nodes[92].neighbors = append(nodes[92].neighbors, nodes[91], nodes[93], nodes[104])
			nodes[93].neighbors = append(nodes[93].neighbors, nodes[77], nodes[92], nodes[105])

			nodes[94].neighbors = append(nodes[94].neighbors, nodes[78], nodes[95], nodes[106])
			nodes[95].neighbors = append(nodes[95].neighbors, nodes[79], nodes[94], nodes[96], nodes[107])
			nodes[96].neighbors = append(nodes[96].neighbors, nodes[80], nodes[95], nodes[108])
			nodes[97].neighbors = append(nodes[97].neighbors, nodes[83], nodes[98], nodes[111])
			nodes[98].neighbors = append(nodes[98].neighbors, nodes[84], nodes[97], nodes[99], nodes[112])
			nodes[99].neighbors = append(nodes[99].neighbors, nodes[85], nodes[98], nodes[113])
			nodes[100].neighbors = append(nodes[100].neighbors, nodes[86], nodes[101], nodes[114])
			nodes[101].neighbors = append(nodes[101].neighbors, nodes[87], nodes[100], nodes[102], nodes[115])
			nodes[102].neighbors = append(nodes[102].neighbors, nodes[88], nodes[101], nodes[116])
			nodes[103].neighbors = append(nodes[103].neighbors, nodes[91], nodes[104], nodes[119])
			nodes[104].neighbors = append(nodes[104].neighbors, nodes[92], nodes[103], nodes[105], nodes[120])
			nodes[105].neighbors = append(nodes[105].neighbors, nodes[93], nodes[104], nodes[121])

			nodes[106].neighbors = append(nodes[106].neighbors, nodes[94], nodes[107], nodes[122])
			nodes[107].neighbors = append(nodes[107].neighbors, nodes[95], nodes[106], nodes[108])
			nodes[108].neighbors = append(nodes[108].neighbors, nodes[96], nodes[107], nodes[109])
			nodes[109].neighbors = append(nodes[109].neighbors, nodes[108], nodes[110])
			nodes[110].neighbors = append(nodes[110].neighbors, nodes[109], nodes[111])
			nodes[111].neighbors = append(nodes[111].neighbors, nodes[97], nodes[110], nodes[112])
			nodes[112].neighbors = append(nodes[112].neighbors, nodes[98], nodes[111], nodes[113])
			nodes[113].neighbors = append(nodes[113].neighbors, nodes[99], nodes[112], nodes[123])
			nodes[114].neighbors = append(nodes[114].neighbors, nodes[100], nodes[115], nodes[124])
			nodes[115].neighbors = append(nodes[115].neighbors, nodes[101], nodes[114], nodes[116])
			nodes[116].neighbors = append(nodes[116].neighbors, nodes[102], nodes[115], nodes[117])
			nodes[117].neighbors = append(nodes[117].neighbors, nodes[116], nodes[118])
			nodes[118].neighbors = append(nodes[118].neighbors, nodes[117], nodes[119])
			nodes[119].neighbors = append(nodes[119].neighbors, nodes[103], nodes[118], nodes[120])
			nodes[120].neighbors = append(nodes[120].neighbors, nodes[104], nodes[119], nodes[121])
			nodes[121].neighbors = append(nodes[121].neighbors, nodes[105], nodes[120], nodes[125])

			nodes[122].neighbors = append(nodes[122].neighbors, nodes[106], nodes[126])
			nodes[123].neighbors = append(nodes[123].neighbors, nodes[113], nodes[127])
			nodes[124].neighbors = append(nodes[124].neighbors, nodes[114], nodes[136])
			nodes[125].neighbors = append(nodes[125].neighbors, nodes[121], nodes[137])

			nodes[126].neighbors = append(nodes[126].neighbors, nodes[122], nodes[138])
			nodes[127].neighbors = append(nodes[127].neighbors, nodes[123], nodes[128], nodes[139])
			nodes[128].neighbors = append(nodes[128].neighbors, nodes[127], nodes[129], nodes[140])
			nodes[129].neighbors = append(nodes[129].neighbors, nodes[128], nodes[130], nodes[141])
			nodes[130].neighbors = append(nodes[130].neighbors, nodes[129], nodes[131], nodes[142], nodes[212])
			nodes[131].neighbors = append(nodes[131].neighbors, nodes[130], nodes[132], nodes[213])
			nodes[132].neighbors = append(nodes[132].neighbors, nodes[131], nodes[133], nodes[214])
			nodes[133].neighbors = append(nodes[133].neighbors, nodes[132], nodes[134], nodes[143], nodes[215])
			nodes[134].neighbors = append(nodes[134].neighbors, nodes[133], nodes[135], nodes[144])
			nodes[135].neighbors = append(nodes[135].neighbors, nodes[134], nodes[136], nodes[145])
			nodes[136].neighbors = append(nodes[136].neighbors, nodes[124], nodes[135], nodes[146])
			nodes[137].neighbors = append(nodes[137].neighbors, nodes[125], nodes[147])

			nodes[138].neighbors = append(nodes[138].neighbors, nodes[126], nodes[148])
			nodes[139].neighbors = append(nodes[139].neighbors, nodes[127], nodes[140], nodes[149])
			nodes[140].neighbors = append(nodes[140].neighbors, nodes[128], nodes[139], nodes[141], nodes[150])
			nodes[141].neighbors = append(nodes[141].neighbors, nodes[129], nodes[140], nodes[142], nodes[151])
			nodes[142].neighbors = append(nodes[142].neighbors, nodes[130], nodes[141], nodes[152])
			nodes[143].neighbors = append(nodes[143].neighbors, nodes[133], nodes[144], nodes[153])
			nodes[144].neighbors = append(nodes[144].neighbors, nodes[134], nodes[143], nodes[145], nodes[154])
			nodes[145].neighbors = append(nodes[145].neighbors, nodes[135], nodes[144], nodes[146], nodes[155])
			nodes[146].neighbors = append(nodes[146].neighbors, nodes[136], nodes[145], nodes[156])
			nodes[147].neighbors = append(nodes[147].neighbors, nodes[137], nodes[157])

			nodes[148].neighbors = append(nodes[148].neighbors, nodes[138], nodes[158])
			nodes[149].neighbors = append(nodes[149].neighbors, nodes[139], nodes[150], nodes[159])
			nodes[150].neighbors = append(nodes[150].neighbors, nodes[140], nodes[149], nodes[151], nodes[160])
			nodes[151].neighbors = append(nodes[151].neighbors, nodes[141], nodes[150], nodes[152], nodes[161])
			nodes[152].neighbors = append(nodes[152].neighbors, nodes[142], nodes[151], nodes[162])
			nodes[153].neighbors = append(nodes[153].neighbors, nodes[143], nodes[154], nodes[163])
			nodes[154].neighbors = append(nodes[154].neighbors, nodes[144], nodes[153], nodes[155], nodes[164])
			nodes[155].neighbors = append(nodes[155].neighbors, nodes[145], nodes[154], nodes[156], nodes[165])
			nodes[156].neighbors = append(nodes[156].neighbors, nodes[146], nodes[155], nodes[166])
			nodes[157].neighbors = append(nodes[157].neighbors, nodes[147], nodes[167])

			nodes[158].neighbors = append(nodes[158].neighbors, nodes[148], nodes[168])
			nodes[159].neighbors = append(nodes[159].neighbors, nodes[149], nodes[160], nodes[169])
			nodes[160].neighbors = append(nodes[160].neighbors, nodes[150], nodes[159], nodes[161], nodes[170])
			nodes[161].neighbors = append(nodes[161].neighbors, nodes[151], nodes[160], nodes[162], nodes[171])
			nodes[162].neighbors = append(nodes[162].neighbors, nodes[152], nodes[161], nodes[172])
			nodes[163].neighbors = append(nodes[163].neighbors, nodes[153], nodes[164], nodes[173])
			nodes[164].neighbors = append(nodes[164].neighbors, nodes[154], nodes[163], nodes[165], nodes[174])
			nodes[165].neighbors = append(nodes[165].neighbors, nodes[155], nodes[164], nodes[166], nodes[175])
			nodes[166].neighbors = append(nodes[166].neighbors, nodes[156], nodes[165], nodes[176])
			nodes[167].neighbors = append(nodes[167].neighbors, nodes[157], nodes[177])

			nodes[168].neighbors = append(nodes[168].neighbors, nodes[158], nodes[178])
			nodes[169].neighbors = append(nodes[169].neighbors, nodes[159], nodes[170], nodes[185])
			nodes[170].neighbors = append(nodes[170].neighbors, nodes[160], nodes[169], nodes[171], nodes[186])
			nodes[171].neighbors = append(nodes[171].neighbors, nodes[161], nodes[170], nodes[172], nodes[187])
			nodes[172].neighbors = append(nodes[172].neighbors, nodes[162], nodes[171], nodes[188])
			nodes[173].neighbors = append(nodes[173].neighbors, nodes[163], nodes[174], nodes[189])
			nodes[174].neighbors = append(nodes[174].neighbors, nodes[164], nodes[173], nodes[175], nodes[190])
			nodes[175].neighbors = append(nodes[175].neighbors, nodes[165], nodes[174], nodes[176], nodes[191])
			nodes[176].neighbors = append(nodes[176].neighbors, nodes[166], nodes[175], nodes[192])
			nodes[177].neighbors = append(nodes[177].neighbors, nodes[167], nodes[199])

			nodes[178].neighbors = append(nodes[178].neighbors, nodes[168], nodes[179])
			nodes[179].neighbors = append(nodes[179].neighbors, nodes[178], nodes[180])
			nodes[180].neighbors = append(nodes[180].neighbors, nodes[179], nodes[181])
			nodes[181].neighbors = append(nodes[181].neighbors, nodes[180], nodes[182])
			nodes[182].neighbors = append(nodes[182].neighbors, nodes[181], nodes[183])
			nodes[183].neighbors = append(nodes[183].neighbors, nodes[182], nodes[184])
			nodes[184].neighbors = append(nodes[184].neighbors, nodes[183], nodes[185])
			nodes[185].neighbors = append(nodes[185].neighbors, nodes[169], nodes[184], nodes[186])
			nodes[186].neighbors = append(nodes[186].neighbors, nodes[170], nodes[185], nodes[187])
			nodes[187].neighbors = append(nodes[187].neighbors, nodes[171], nodes[185], nodes[188])
			nodes[188].neighbors = append(nodes[188].neighbors, nodes[172], nodes[187])
			nodes[189].neighbors = append(nodes[189].neighbors, nodes[173], nodes[190])
			nodes[190].neighbors = append(nodes[190].neighbors, nodes[174], nodes[189], nodes[191])
			nodes[191].neighbors = append(nodes[191].neighbors, nodes[175], nodes[190], nodes[192])
			nodes[192].neighbors = append(nodes[192].neighbors, nodes[176], nodes[191], nodes[193])
			nodes[193].neighbors = append(nodes[193].neighbors, nodes[192], nodes[194])
			nodes[194].neighbors = append(nodes[194].neighbors, nodes[193], nodes[195])
			nodes[195].neighbors = append(nodes[195].neighbors, nodes[194], nodes[196])
			nodes[196].neighbors = append(nodes[196].neighbors, nodes[195], nodes[197])
			nodes[197].neighbors = append(nodes[197].neighbors, nodes[196], nodes[198])
			nodes[198].neighbors = append(nodes[198].neighbors, nodes[197], nodes[199])
			nodes[199].neighbors = append(nodes[199].neighbors, nodes[177], nodes[198])

			nodes[200].neighbors = append(nodes[200].neighbors, nodes[201], nodes[204])
			nodes[201].neighbors = append(nodes[201].neighbors, nodes[200], nodes[202], nodes[205])
			nodes[202].neighbors = append(nodes[202].neighbors, nodes[201], nodes[203], nodes[206])
			nodes[203].neighbors = append(nodes[203].neighbors, nodes[202], nodes[207])

			nodes[204].neighbors = append(nodes[204].neighbors, nodes[200], nodes[205], nodes[208])
			nodes[205].neighbors = append(nodes[205].neighbors, nodes[201], nodes[204], nodes[206], nodes[209])
			nodes[206].neighbors = append(nodes[206].neighbors, nodes[202], nodes[205], nodes[207], nodes[210])
			nodes[207].neighbors = append(nodes[207].neighbors, nodes[203], nodes[206], nodes[211])

			nodes[208].neighbors = append(nodes[208].neighbors, nodes[204], nodes[209], nodes[212])
			nodes[209].neighbors = append(nodes[209].neighbors, nodes[205], nodes[208], nodes[210], nodes[213])
			nodes[210].neighbors = append(nodes[210].neighbors, nodes[206], nodes[209], nodes[211], nodes[214])
			nodes[211].neighbors = append(nodes[211].neighbors, nodes[207], nodes[210], nodes[215])

			nodes[212].neighbors = append(nodes[212].neighbors, nodes[208], nodes[213], nodes[130])
			nodes[213].neighbors = append(nodes[213].neighbors, nodes[209], nodes[212], nodes[214], nodes[131])
			nodes[214].neighbors = append(nodes[214].neighbors, nodes[210], nodes[213], nodes[215], nodes[132])
			nodes[215].neighbors = append(nodes[215].neighbors, nodes[211], nodes[214], nodes[133])

		}

		imd.Color = colornames.Yellow
		for _, singlenode := range nodes {
			if !singlenode.isEaten {
				imd.Push(singlenode.vec)
			}
		}
		imd.Circle(7, 0)

	}
}

func squareRun(currentNodeGhost *Node, win *pixelgl.Window) {
	imd.Color = colornames.Red

	timeToInstruction := 80.5
	timeToJump := 0.16
	last := time.Now()
	sumTimePassed := timeToInstruction
	sumTimePassedJump := timeToJump
	posInPath := 0
	var ghostPath []*Node
	pos := currentNodeGhost.vec

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		sumTimePassed += dt
		sumTimePassedJump += dt
		checkCollisionGhost(pos)
		if sumTimePassed >= timeToInstruction {
			posInPath = 0
			for _, singleNode := range nodes {
				if checkCollision(pos, singleNode.vec) {
					currentNodeGhost = singleNode
					break
				}
			}
			if currentNodeGhost != nil {
				ghostPath = Breadthwise(*currentNodeGhost, *currentNodePacman)
			}
			sumTimePassed = 0.0
		}

		if sumTimePassedJump >= timeToJump {
			if len(ghostPath) > 1 {
				currentNodeGhost = ghostPath[posInPath]
				pos = ghostPath[posInPath].vec
				posInPath++
				sumTimePassedJump = 0
				if posInPath >= len(ghostPath) {
					sumTimePassed += 100
					posInPath = 0
				}
			} else {

			}

		}

		square := imdraw.New(nil)
		square.Color = colornames.Red
		rect := pixel.R(pos.X, pos.Y, pos.X+20, pos.Y+20)
		square.Push(rect.Min, rect.Max)
		square.Rectangle(0)
		square.Draw(win)
	}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pacman",
		Bounds: pixel.R(0, 0, 1024, 720),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	setUpLevel()
	currentNodePacman = nodes[20]

	pacmanDown, err := loadPicture("PacmanDown.png")
	if err != nil {
		panic(err)
	}
	pacmanUp, err := loadPicture("PacmanUp.png")
	if err != nil {
		panic(err)
	}
	pacmanLeft, err := loadPicture("PacmanLeft.png")
	if err != nil {
		panic(err)
	}
	pacmanRight, err := loadPicture("PacmanRight.png")
	if err != nil {
		panic(err)
	}

	pacmanSpriteUp := pixel.NewSprite(pacmanUp, pacmanUp.Bounds())
	pacmanSpriteDown := pixel.NewSprite(pacmanDown, pacmanDown.Bounds())
	pacmanSpriteLeft := pixel.NewSprite(pacmanLeft, pacmanLeft.Bounds())
	pacmanSpriteRight := pixel.NewSprite(pacmanRight, pacmanRight.Bounds())

	currentPos = pixel.V(350, 40) //Starting position

	shouldMoveLeft := true
	shouldMoveRight := false
	shouldMoveUp := false
	shouldMoveDown := false

	hasWallLeft := false
	hasWallRight := false
	hasWallUp := false
	hasWallDown := false

	for i := 0; i < numGhosts; i++ {
		go squareRun(nodes[(200+i)%215], win)
	}

	nodesEaten := 0
	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		win.Clear(colornames.Black)

		mat := pixel.IM

		for _, singlenode := range nodes { //Draw and check collision with node
			if checkCollision(currentPos, singlenode.vec) {
				currentNodePacman = singlenode
				if !singlenode.isEaten {
					imd.Color = colornames.Black
					imd.Push(singlenode.vec)
					imd.Circle(8, 0)
					score += 100
					nodesEaten++
					fmt.Printf("Score: %d \n", score)
					if nodesEaten == len(nodes)-16 {
						log.Println("YOU WIN!!")
						os.Exit(0)
					}
				}
				singlenode.isEaten = true
			}
		}
		mat = mat.ScaledXY(pixel.ZV, pixel.V(0.72, 0.72))
		imd.Draw(win)

		if shouldMoveLeft {
			if !pause && currentPos.X > 25 && !hasWallLeft {
				currentPos.X -= speed * float64(dt)
			}
			mat = mat.Moved(currentPos)
			pacmanSpriteLeft.Draw(win, mat)
			if checkCollisionWall(currentPos, pixel.R(0, 360-30, 20, 360+30)) {
				currentPos.X = 1024
			}
			for _, singeWall := range walls {
				if checkCollisionWall(currentPos, singeWall.rect) && !hasWallRight && !hasWallUp && !hasWallDown {
					hasWallLeft = true
				}
			}
			if !checkCollisionAllWalls(currentPos) {
				hasWallUp = false
				hasWallRight = false
				hasWallDown = false
			}
		} else if shouldMoveRight {
			if !pause && currentPos.X < 1000 && !hasWallRight {
				currentPos.X += speed * float64(dt)
			}
			mat = mat.Moved(currentPos)
			pacmanSpriteRight.Draw(win, mat)
			if checkCollisionWall(currentPos, pixel.R(1000, 360-30, 1024, 360+30)) {
				currentPos.X = 0
			}
			for _, singeWall := range walls {
				if checkCollisionWall(currentPos, singeWall.rect) && !hasWallLeft && !hasWallUp && !hasWallDown {
					hasWallRight = true
				}
			}
			if !checkCollisionAllWalls(currentPos) {
				hasWallLeft = false
				hasWallUp = false
				hasWallDown = false
			}
		} else if shouldMoveUp {
			if !pause && currentPos.Y < 695 && !hasWallUp {
				currentPos.Y += speed * float64(dt)
			}
			mat = mat.Moved(currentPos)
			pacmanSpriteUp.Draw(win, mat)
			for _, singeWall := range walls {
				if checkCollisionWall(currentPos, singeWall.rect) && !hasWallRight && !hasWallLeft && !hasWallDown {
					hasWallUp = true
				}
			}
			if !checkCollisionAllWalls(currentPos) {
				hasWallLeft = false
				hasWallRight = false
				hasWallDown = false
			}
		} else if shouldMoveDown {
			if !pause && currentPos.Y > 25 && !hasWallDown {
				currentPos.Y -= speed * float64(dt)
			}
			mat = mat.Moved(currentPos)
			pacmanSpriteDown.Draw(win, mat)
			for _, singeWall := range walls {
				if checkCollisionWall(currentPos, singeWall.rect) && !hasWallRight && !hasWallUp && !hasWallLeft {
					hasWallDown = true
				}
			}
			if !checkCollisionAllWalls(currentPos) {
				hasWallLeft = false
				hasWallRight = false
				hasWallUp = false
			}
		}

		if win.Pressed(pixelgl.KeyLeft) {
			if pause {
				currentPos.X -= speed
			}
			pause = false
			shouldMoveLeft = true
			shouldMoveRight = false
			shouldMoveUp = false
			shouldMoveDown = false
			hasWallRight = false
		} else if win.Pressed(pixelgl.KeyRight) {
			if pause {
				currentPos.X += speed
			}
			pause = false
			shouldMoveLeft = false
			shouldMoveRight = true
			shouldMoveUp = false
			shouldMoveDown = false
			hasWallLeft = false
		} else if win.Pressed(pixelgl.KeyDown) {
			if pause {
				currentPos.Y -= speed
			}
			pause = false
			shouldMoveLeft = false
			shouldMoveRight = false
			shouldMoveUp = false
			shouldMoveDown = true
			hasWallUp = false
		} else if win.Pressed(pixelgl.KeyUp) {
			if pause {
				currentPos.Y += speed
			}
			pause = false
			shouldMoveLeft = false
			shouldMoveRight = false
			shouldMoveUp = true
			shouldMoveDown = false
			hasWallDown = false
			if !checkCollisionAllWalls(currentPos) {
				hasWallLeft = false
				hasWallRight = false
			}

		}
		win.Update()
	}
}

func (n Node) compare(m Node) bool {
	if n.id == m.id {
		return true
	}
	return false
}

func contains(s []*Node, e *Node) bool {
	for _, a := range s {
		if a.compare(*e) {
			return true
		}
	}
	return false
}

func Breadthwise(start, end Node) []*Node {
	size := 300

	start.history = start.history[:0]
	result := make([]*Node, 0, size)

	visited := make([]*Node, 0, size)

	work := make([]*Node, 0, size)

	visited = append(visited, &start)       //visited.Add(start)
	work = append([]*Node{&start}, work...) //work.Enqueue(start)
	for len(work) > 0 {
		current := work[len(work)-1]
		work = work[:len(work)-1] //current = work.Dequeue
		if current.compare(end) {
			//Found node
			result = current.history
			result = append(result, current)
			return result
		}
		//Didnt find node
		for i := 0; i < len(current.neighbors); i++ {
			currentNeighbor := current.neighbors[i]
			if !contains(visited, currentNeighbor) {
				currentNeighbor.history = make([]*Node, len(current.history), size)
				copy(currentNeighbor.history, current.history)
				currentNeighbor.history = append(currentNeighbor.history, current)
				visited = append(visited, currentNeighbor)
				work = append([]*Node{currentNeighbor}, work...)
			}
		}
	}
	return nil
}

func Depthwise(start, end Node) []*Node {
	size := 300

	result := make([]*Node, 0, size)
	visited := make([]*Node, 0, size)
	work := make([]*Node, 0, size)
	var current *Node
	var currentSon *Node

	start.history = start.history[:0]
	visited = append(visited, &start)
	work = append([]*Node{&start}, work...) //work.Enqueue(start)

	for len(work) > 0 {
		current, work = work[len(work)-1], work[:len(work)-1]
		if current.compare(end) {
			result = current.history
			result = append(result, current)
			return result
		} else {
			for _, node := range current.neighbors {
				currentSon = node
				if !contains(visited, currentSon) {
					visited = append(visited, currentSon)
					currentSon.history = current.history
					currentSon.history = append(currentSon.history, current)
					work = append([]*Node{currentSon}, work...) //work.push(currentSon)
				}
			}
		}
	}
	return nil
}

func main() {
	//usage: go run pacman.go -g <numGhosts>
	numGhosts, _ = strconv.Atoi(os.Args[2])
	rand.Seed(time.Now().UnixNano())
	pixelgl.Run(run)

	/*setUpLevel()

	res := Depthwise(*nodes[0], *nodes[20])
	for _, n := range res {
		fmt.Printf("-> %d ", n.id)
	}
	fmt.Println()

	res = Breadthwise(*nodes[0], *nodes[20])
	for _, n := range res {
		fmt.Printf("-> %d ", n.id)
	}
	fmt.Println()*/
}
