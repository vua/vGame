package component

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math"
	"math/rand"
	"strconv"
)

type Square struct {
	bgc   color.RGBA
	h     float64
	w     float64
	x     float64
	y     float64
	step  float64
	angle float64
	stepX float64
	stepY float64
	score int
	alive bool
	IsRun bool
	Image *ebiten.Image
	Opts  *ebiten.DrawImageOptions
}

func NewSquare(bgc color.RGBA, h, w int, x, y, step float64) *Square {
	image := ebiten.NewImage(w, h)
	image.Fill(bgc)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(x, y)
	angle := float64(rand.Intn(120) + 30)
	return &Square{
		bgc:   bgc,
		h:     float64(h),
		w:     float64(w),
		x:     x,
		y:     y,
		step:  step,
		angle: angle,
		stepX: step * math.Cos(angle),
		stepY: -step * math.Sin(angle),
		alive: true,
		Image: image,
		Opts:  opts,
	}
}
func (s *Square) GetScore() string{
	return strconv.Itoa(s.score)
}

func (s *Square) IsAlive() bool{
	return s.alive
}

func (s *Square) CollisionDetection(recv *Square,w,h float64) {
	x, y := s.x+s.stepX, s.y+s.stepY
	tx,ty:=s.stepX,s.stepY
	//碰撞垂直边 (-stepX,stepY)
	//碰撞水平边 (stepX,-stepY)
	if x<=0 {
		tx=-s.x
		ty=tx*math.Tan(s.angle)
		s.stepX*=-1
	} else if x+s.w>=w{
		tx=320-s.w-s.x
		ty=tx*math.Tan(s.angle)
		s.stepX*=-1
	}else if y<=0 {
		ty=-s.y
		tx=ty/math.Tan(s.angle)
		s.stepY*=-1
	}else if y+s.h >= recv.y&& y<= recv.y+recv.h&& x <= recv.w+recv.x && x+s.w >= recv.x {
		s.score++
		ty=recv.y-s.y-s.w
		tx=ty/math.Tan(s.angle)
		if s.stepY*s.stepX<0 {
			tx*=-1
		}
		s.stepY*=-1
	}else if y+s.h>=h {
		s.alive=false
		return
	}
	s.x+=tx
	s.y+=ty
	s.Opts.GeoM.Translate(tx, ty)
}

func (s *Square) HitDetection(awards *[]*Square){
	for i:=0;i<len(*awards);i++{
		award:=(*awards)[i]
		if math.Abs(s.x-award.x)<=s.w&&math.Abs(s.y-award.y)<=s.h {
			*awards=append((*awards)[0:i],(*awards)[i+1:]...)
			i--
			s.score+=10
		}
	}
}

func (s *Square) Move(w int,step float64,boll *Square) {
	W:=float64(w)

	if s.x+step < 0 {
		s.Opts.GeoM.Translate(0-s.x, 0)
		if !boll.IsRun{
			boll.Opts.GeoM.Translate(0-s.x, 0)
			boll.x+=-s.x
		}
		s.x = 0
		return
	}
	if s.x+step > W-s.w {
		s.Opts.GeoM.Translate(W-s.w-s.x, 0)
		if !boll.IsRun{
			boll.Opts.GeoM.Translate(W-s.w-s.x, 0)
			boll.x+=W-s.w-s.x
		}
		s.x = W - s.w
		return
	}
	s.x += step
	s.Opts.GeoM.Translate(step, 0)
	if !boll.IsRun {
		boll.Opts.GeoM.Translate(step, 0)
		boll.x+=step
	}
}
