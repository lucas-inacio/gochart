# gochart
Simple chart implementation using fyne and gg libraries

### Dependencies
* gg (https://github.com/fogleman/gg)
* fyne (https://github.com/fyne-io/fyne)

### Functionality
Currently only BarChart and LineChart are provided. More types can be created by embedding BaseChart and implementing the required methods:
* Required by Chart interface
```go
      // A "contructor" like NewAwesomeChart should be created as well
      BeginDraw(dc *gg.Context, x, y float64)
      DrawPoint(dc *gg.Context, x, y, yOffset float64, int, image.Rectangle)
      EndDraw(dc *gg.Context)
 ```
    
 * Required by Fyne (follow the guidelines at [Extending Widgets](https://developer.fyne.io/tutorial/extending-widgets)) for more information:
```go
      CreateRenderer() fyne.WidgetRenderer
 ```
