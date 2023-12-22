---
title: "Building a fintech Banking application with React(typescript) and Golang(GIN)"
description: "This is a series where I will be building a banking application from scratch using the React libaray with typescript, Golang,Gin postgres sqlite and redis "
image: "images/post/go-gin-api.jpg"
date: 2023-06-24T18:19:25+06:00
categories: ["programming", "Reguler Club"]
tags: ["go", "gin", "Vue", "Gorm"]
type: "featured" # available types: [featured/regular]
draft: false
---


Hello guys, welcome to the start of a series where I will be documenting, how to use Go(GIN), React(typesscript) to build a banking application(Fintech). I started this project to learn more about using react and typescript, aswell as to improve my go knowledge, so I am learning as I am building. I'm hoping this tutorial will serve useful to you and you learn somethiing with me. Anyways, sit tight, grab your keyboard and let's code.

#### Project breakdown
A little bit about the project we will be building. It's a simple banking application, the bankend will be built in Go with a monolith architecture, although I do plan on exploring seperating into micro services once I'm done building this application, if you're interested in learning that with me then don't miss out and give me a follow. I will try my best to explauin concepts as much as I can, so even a begineer can follow, but you will need a bit of GO knowlege and Javascript knowlege to follow up.        


The name of this project is Reguler Club, you may have seen or heard about the cartoon network series `reguler show`

{{< image src="https://deadline.com/wp-content/uploads/2016/09/regular-show.jpg?w=681&h=383&crop=1" caption="Reguler show" alt="alter-text" height="" width="" position="center" command="fill" option="q100" class="img-fluid" title="modecai and rigby" webp="false" >}}

I absolutely loved that cartoon, the idea of this project was somewhat gotten from that. Reguler Club is like a community forum for `nerds` who like all sorts of cartoons, anime, comics, video games or computers to come hangout and share ideas and discussions. But since, the project was made for learning purposes I plan on integrating certain different features, like working with transactions, in our case we will create our own little fake money called `Gokens`, we will create a whole transactions system, where users can buy gokens, don't worry its free, all they have to do is fill a form with the almighty phrase `GO IS COOL`, and they get to claim a free `goken`, but only once every 24 hrs. Gokens have a heirachy in values and rarity, they can be traded among users and randomly increase and decrease in value based on our algorithm we will implement and the how much its being traded. Gokens can also be shared among users, and can be given as gifts to other users, Gokens can also be used to unlock certain Locked features in our application. Which brings us to another feature, we will different roles assigned to users, an admin, mod, prime(high value gokens), and normal authenticated users, just like in many real world applications. Different roles will give the user different access to what they can do on our application.
We will also have real-time messaging and group communication. It's a simple project overall but there's a lot to learn from it. I hope you guys are excited to partake on this series with me. 

#### To begin here's a bit of things that you'll need.
First off, I'm assuming you know a atleast a tiny little bit of GO and javascript, or any other programming language knowlege will be useful. I will try to break steps down as I possibly can do so feel free to join me regardless.

Here's what you'll need:
A computer `obviously`, you're not gonna write code on a rock, OR can you?? hmmmm. Your OS sholdnt matter so far you're on a computer, I heard you could also code on an Andriod with [Termux](https://termux.dev/en/),
[read](https://github.com/golang/go/wiki/Mobile), if you don't have a laptop or desktop. If you're on ios, I'm sorry, I don't know.

If you do have a laptop/desktop:
You can download and install Go from the [official website](https://golang.org/).
We'll be using Vue for the front end so we'll need [Node](https://nodejs.org/en/download) installed aswell.

You'll also need to install [docker](https://www.docker.com/get-started). desktop although its not entirely necessary, for now we will be only using docker only for our database (postgres). But you're free to use whatever database you want or whatever way you want to host it, so far you know how to set it up correctly. I won't be explaining that here but if I ever make a tutorial on that in the future, I will link it up here. 
The gorm [documentation](https://gorm.io/docs/connecting_to_the_database.html) should be helpful for setting up a relational DB with gorm(the go to ORM package in golang).

We'll be installing a bunch of other stuff in the future but for now let's keep it simple. 

#### check things out
If you've succesfully installed go you should be able to run this simple program. First, let's open or terminal, input the following command to create a new directory where we will run our simple hello world, just to test things out.

{{< highlight bash "linenos=table,hl_lines=8 15-17,linenostart=22" >}}
   # Create a new repository
mkdir helloword
cd helloword
code

# Initialize the project with Go modules
go mod init github.com/yourcoolgithubname/repo

{{< /highlight >}}


#### Installing Dependencies

Now, let's install the necessary dependencies, including the Gin framework:
{{< highlight bash "linenos=table,hl_lines=8 15-17,linenostart=22" >}}
go get -u github.com/gin-gonic/gin
{{< /highlight >}}

#### print Hello world

Let's create a main file in our root directly, in Go, everything is ran from the main.go file. `main.go`:

{{< highlight go "linenos=table,hl_lines=8 15-17,linenostart=22" >}}package main

import (
	"fmt"
)
func main() {
    message := "Hello world"
    fmt.Println(message)
}

{{< /highlight >}}

If that outputs `Hello world` successfully, then you're good to `Go` with us. 

Ok that has been enough introduction, in this next chapter we will be diving into the real stuff...
[Building a fullstack application with Go Gin Gorm and Vue js, setting up Postgres and models 1](http://charlesdpj.com/blog/Gin-Gorm-Golang-Vue-jwt-authetication.md)