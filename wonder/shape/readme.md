


```
plane := shape.Plane{}

a, b := f32.Point(0, 100)}, f32.Point{0, 100}

line := shape.Line{a, b}

plane.Add(line)

text := shape.Text{}

plane.Add(text)

// 
plane.View(f32.Rect(0, 0, 400, 600), gtx)

```

Zoom in out 

```

```

Pan though plane

```
offset := f32.Point{20, 30}
viewport.Add(offset)
```

Select Shape

Draw Polyline

## Research

* [Draw sin wave with brezier curves](https://stackoverflow.com/questions/29022438/how-to-approximate-a-half-cosine-curve-with-bezier-paths-in-svg)
