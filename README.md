# Menu Scraper
This _scraper_ get launch menu from this three restaurants from Brno: https://www.pivnice-ucapa.cz/denni-menu.php, https://www.suzies.cz/poledni-menu and https://www.menicka.cz/4921-veroni-coffee--chocolate.html .

Program is simple with comments. I have one issue with interpreter. I must set `export GO111MODULE="auto"` to resolve the error with modules.

I use goquery to scrap the HTML document.

## Instruction
Get the 3rd party libs:
```
go get golang.org/x/text/encoding/charmap
go get github.com/PuerkitoBio/goquery
```
Then you can run the program 
```
go run main.go
```