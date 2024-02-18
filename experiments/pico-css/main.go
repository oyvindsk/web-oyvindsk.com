package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {

	html := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta name="generator" content="HTML Tidy for HTML5 for Linux version 5.7.45">
  <title>Øyvind Gerrard Skaar || oyvindsk.com</title>
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="stylesheet" href="mystyle.css">
</head>
<body>
  <article>
    <header>
      <nav>
        <a href="/" title="Home">Øyvind G. Skaar</a> <a href="/writing" title="Articles">Articles</a> <a href="/hire-me" title="Hire Me">Hire Me</a> <a href="/projects" title="Projects &amp; Clients">Projects & Clients</a> <a href="/contact" title="Contact">Contact</a>
      </nav>
    </header>

	<p>PPPP</p>
	<p>PPPP2 </p>

  </article>


	
<!-- Code injected by live-server -->
<script type="text/javascript">
	// <![CDATA[  <-- For SVG support
	if ('WebSocket' in window) {
		(function() {
			function refreshCSS() {
				var sheets = [].slice.call(document.getElementsByTagName("link"));
				var head = document.getElementsByTagName("head")[0];
				for (var i = 0; i < sheets.length; ++i) {
					var elem = sheets[i];
					head.removeChild(elem);
					var rel = elem.rel;
					if (elem.href && typeof rel != "string" || rel.length == 0 || rel.toLowerCase() == "stylesheet") {
						var url = elem.href.replace(/(&|\?)_cacheOverride=\d+/, '');
						elem.href = url + (url.indexOf('?') >= 0 ? '&' : '?') + '_cacheOverride=' + (new Date().valueOf());
					}
					head.appendChild(elem);
				}
			}
			var protocol = window.location.protocol === 'http:' ? 'ws://' : 'wss://';
			var address = protocol + window.location.host + window.location.pathname + '/ws';
			var socket = new WebSocket(address);
			socket.onmessage = function(msg) {
				if (msg.data == 'reload') window.location.reload();
				else if (msg.data == 'refreshcss') refreshCSS();
			};
			console.log('Live reload enabled.');
		})();
	}
	// ]]>
</script>
</body>
</html>
`

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, html)
	})

	log.Println("Staring on port 8080")
	e.Logger.Fatal(e.Start(":1323"))

}
