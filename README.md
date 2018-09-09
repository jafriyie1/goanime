# goanime
My efforts to provide a quick and easy anime viewing experience 

## Motivation
I absolutely love watching anime, and I have been using a plethora of different anime sites for years. But recently 
I have been getting rather annoyed about all of the intrusive/ adult content ads, unresponsive webpages, and waiting 
for a long time to just watch a show. Sure I could write scripts so that I can actually use 
Adblock on the sites, but that honestly is a hassell. What if the viewing experience could be quick, easy and painless?

This repo contains code that provides a quick and easy viewing experience. 
Currently I am using chromedp in headless chrome mode, and other packages to get the searched show. 
The code then opens up a new tab in your preferred browser with the selected show. It bypasses any ads 
and gets the show that you want to watch. 

## Contributors 
If you would like to contribute that would be absolutely swell :) !!! I am currently learning Go, as 
I think it is a beautiful language; because of this, the code may not be written in the best way. With this in mind
anyone who has ideas for features or how to improve the code would be greatly appreciated!

## Example of program in terminal (with logging) 
[![asciicast](https://asciinema.org/a/KttZeSMSQ2musQVoPh2lr8MDI.png)](https://asciinema.org/a/KttZeSMSQ2musQVoPh2lr8MDI)

## How to use 
To use the streamer, clone the repo. If you have a mac, navigate to the bin folder and run 
the binary there (execute ./goanime_concurrent in your terminal, or double click it). If you have a Windows machine run the .exe file in the bin folder.

## Disclaimer 

As of right now this application is still very raw. I can't guarantee that it will work for your machine due to the fact that I haven't written full tests for potential pitfalls. This program also only works for subbed anime. That will change soon.

## Upcoming Features
- Allow Dubbed versions of shows (Not Completed)
- Allow users to put more than one episode (Finished! But it is still in Alpha!)
- Include search feature of possible shows (Finished!)
- Create GUI App
- Easy installation and cross compatability (Finished!)
