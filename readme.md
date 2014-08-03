# JSON to XML

This is a simple go web application that uses `github.com/clbanning/mxj` and `github.com/clbanning/anyxml` package to convert JSON to XML and visa versa. This is ready to be deployed on heroku.

I created this because I was working with some web services in PeopleSoft and the JSON support was terrible. Peoplesoft could generate valid JSON but it was having problems parsing.  I figured I could get the JSON from the web service. Then invoke this service to get XML back which could then more easily be parsed by PeoplelSoft.

Special thanks to @clbanning for both of his libraries and help.


## Issue List

* needs some more test cases around XML output
** Thinking that the text cases should have expect output slice of xml nodes and values since you cannot expect an order to come back so you cannot compare full XML string. 




## Deploying on Heroku

1. First see [getting started with go on heroku](http://mmcgrana.github.io/2012/09/getting-started-with-go-on-heroku.html)


`heroku create -b https://github.com/kr/heroku-buildpack-go.git`
`git push heroku master`