


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
* [Check if point is inside of a rectangle](https://math.stackexchange.com/questions/190111/how-to-check-if-a-point-is-inside-a-rectangle)
* [Collision Detection](http://www.jeffreythompson.org/collision-detection/table_of_contents.php)
* [JuliaMono](https://juliamono.netlify.app/), a monospaced typeface designed for programming in Julia, [HN](https://news.ycombinator.com/item?id=24732973) comments
