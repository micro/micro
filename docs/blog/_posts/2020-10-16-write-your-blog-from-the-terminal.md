---
layout:	post
author: Janos Dobronszki
title:	"Write your blog from the terminal"
date:	2020-11-09 10:00:00
---

If you're following our [Building a Blog with Micro](/blog/2020/11/06/building-a-blog-with-micro-part-one) series you'll already know a lot 
about building a blog as Micro services. In this post we're going to cover a little on actually writing your blog via the terminal. We've already 
built out the services needed so you should be able to access them in [micro/services/blog](https://github.com/micro/services/tree/master/blog).

## Deploying the services

Make sure your env is set to dev and deploy the following services:

```sh
$ micro env set dev
# micro signup if you don't have an account
$ micro env
  local      127.0.0.1:8081
* dev        proxy.m3o.dev
  platform   proxy.m3o.com
$ micro run github.com/micro/services/blogs/tags
$ micro run github.com/micro/services/blogs/posts
```

If you don't have an account, sign up to the dev env for a free account with the `micro signup` command after setting the env with `micro env set dev`.

After `micro status` confirms that the services are running, the services will be ready to be interacted with.

## Terminal usage

Saving a post:
```sh
$ micro posts save --id=1 --title="Awesome post title" --content="Even more awesome post content"
```

Listing posts
```
$ micro posts query
```

Listing tags
```
$ micro tags list
```

For a comprehensive list of commands, please refer to the readmes in the [`posts`](https://github.com/micro/services/tree/master/blog/posts) and [`tags`](https://github.com/micro/services/tree/master/blog/tags) folders.

## Building a Netlify frontend on the Micro backend

To serve as an example, we have cooked up a very minimal (consisting of a few dozen lines of code) Angular application that uses some endpoints of these services.
The web app in action can be found [here](https://loving-goodall-44ee08.netlify.app/), and the source code for it can be found [here](https://github.com/micro/dev/tree/master/blog/frontend).

You can easily deploy on Netlify with the following settings:

![Netlify settings](/blog/images/netlify-build-settings.png)

Although example code is an Angular one, since it's using the JSON HTTP API, it's trivial to call the backend in any framework in language. The following curl contacts the posts service running in the author's namespace:

```sh
$ curl -H "Micro-Namespace: hatchling-mutation-alright" https://api.m3o.dev/posts/query
{
  "posts": [
    {
      "id": "1",
      "title": "Welcome to my Micro blog hosted on M3o",
      "slug": "welcome-to-my-micro-blog-hosted-on-m3o",
      "content": "The Dev environment of M3o is free to use and a perfect place to host your blog. All it takes is a 'micro run blog/posts && micro run blog/tags' command and you are ready to write blog posts from the terminal. In fact this entry was written in the terminal too."
    },
    {
      "id": "2",
      "title": "Build a blog on our headless CMS",
      "slug": "build-a-blog-on-our-headless-cms",
      "content": "To deploy your blog all you need is a place that can host HTML files and a free M3o Dev account. Since HTML hosting is free on Github and Netlify, you can, in minutes, host your fancy blog that can be managed through the terminal."
    }
  ]
}
```

Happy hacking!

The Micro Team

