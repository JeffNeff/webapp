# How I Create Applications With a React Frontend and Golang Backend
I wanted to share with the community how I create applications with a React frontend and Golang backend. This readme contains the steps to reproduce this repository.

## Development Prerequisites
* Go
* React
* create-react-app
* react-scripts

## Building/Deploying Prerequisites
* ko

## Initalize the project
Steps:
* Create a new folder in your $GOPATH and name it `webapp`.
* Run `go mod int && go mod tidy`.
* Create the following folder/file structure.
```
webapp/
├─ backend/
│  ├─ cmd/
│  │  ├─ webapp/
│  │  │  ├─ main.go
│  ├─ pkg/
│  │  ├─ controller.go
├─ frontend/
├─ Dockerfile
├─ Makefile
```
* Add dependencies

This can be accomplished by executing the following commands:

```
mkdir webapp
cd webapp
mkdir -p backend/cmd/webapp
mkdir -p backend/pkg/
touch Makefile
touch backend/cmd/webapp/main.go
touch backend/pkg/controller.go
touch backend/Dockerfile
go mod init
go mod tidy
npx create-react-app frontend
cd frontend
yarn add react-scripts
```

## Populate the Go files

We need to create a simple Go server to serve our React Frontend.

Navigate to the `backend/cmd/webapp/main.go` file and update the contents to the following:
```
package main

import (
	"log"
	"net/http"

    // this will need to be updated to reflect where your controller is located
	controller "github.com/JeffNeff/webapp/backend/pkg"
)

func main() {

	c := &controller.Controller{}

	http.HandleFunc("/", c.RootHandler)

	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
```

This route will be defined in the `backend/pkg/controller.go` file. Update this file to the following:

```
package controller

import (
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

type Controller struct {
	rootHandler http.Handler
}

var once sync.Once

func (c *Controller) RootHandler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		kdp := path.Join(os.Getenv("KO_DATA_PATH"))
		if !strings.HasSuffix(kdp, "/") {
			kdp = kdp + "/"
		}
		c.rootHandler = http.FileServer(http.Dir(kdp))
	})
	c.rootHandler.ServeHTTP(w, r)
}
```

## Populate the Makefile
I dont like to remember long commands, so often I will populate a Makefile with the commands to build and run the application.

Populate the `Makefile` with the following:
```
update-static:
	@cd backend/cmd/webapp && rm -rf kodata
	@cd frontend && yarn install && yarn react-scripts build && mv build kodata && mv kodata ../backend/cmd/webapp/kodata

run:
	@KO_DATA_PATH=backend/cmd/webapp/kodata/ go run backend/cmd/webapp/main.go

debug:
	@make update-static
	@make run
```

## Build the Static Site
Now we need to build the static site.
```
make update-static
```

## Runing Locally
Lets see what we have so far! You can run the application locally with:
```
make run
```
and then navigate to `http://localhost:8080/` to find the React app we created in the previous step.


# Sending requests to the backend
Now we need to send requests to the backend. Lets modify

## Update the Frontend to Accept Input
First lets add some more dependecys to the frontend project:
```
cd frontend
yarn add @material-ui/core
yarn add axios
```

Now lets create a `PostForm.js` file in the `frontend/src` directory to hold the form we will use to send data to our backend.

```
touch frontend/src/PostForm.js
```

Populate this file with the following:

```
import React from "react";
import { TextField, FormControl, Button } from "@material-ui/core";
import axios from "axios";

const PostForm = () => {
  const [input, setInput] = React.useState("");
  const [response, setResponse] = React.useState("");
  const sendData = (event) => {
    axios
      .post("/event", {
        input,
      })
      .then(function (response) {
        console.log(response);
        if (response.data !== null) {
          setResponse(response.data);
        }
      })
      .catch(function (error) {
        console.log(error);
      });
  };

  return (
    <div>
      <FormControl>
        <TextField
          id="dummy-input"
          label="Your Name"
          value={input}
          onChange={(e) => setInput(e.target.value)}
        />
        <Button onClick={sendData}>Send</Button>
      </FormControl>
      <div>Response: {response}</div>
    </div>
  );
};

export default PostForm;
```

Now lets use this form by modifying the `App.js` file in the `frontend/src` directory to the following:

```
import './App.css';
import PostForm from './PostForm';

function App() {
  return (
    <div className="App">
      <PostForm />
    </div>
  );
}

export default App;
```

## Update the Backend to Accept Input
Now that our frontend is able to post data to the `/event` route, we need to update the backend to accept this data and respond with a message.

Update the `backend/cmd/webapp/main.go` file to the following to handle the `/event` route:

```
package main

import (
	"log"
	"net/http"

	controller "github.com/JeffNeff/webapp/backend/pkg"
)

func main() {
	c := &controller.Controller{}

	http.HandleFunc("/", c.RootHandler)
	http.HandleFunc("/event", c.HandlePost)

	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
```

Update the `backend/pkg/controller.go` file to the following to include this route:
```
package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

type Controller struct {
	rootHandler http.Handler
}

type event struct {
	Input string `json:"input"`
}

var once sync.Once

func (c *Controller) RootHandler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		kdp := path.Join(os.Getenv("KO_DATA_PATH"))
		if !strings.HasSuffix(kdp, "/") {
			kdp = kdp + "/"
		}
		c.rootHandler = http.FileServer(http.Dir(kdp))
	})
	c.rootHandler.ServeHTTP(w, r)
}

func (c *Controller) HandlePost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error occured reading body: %v", err)
		json.NewEncoder(w).Encode("Error reding request")
		return
	}

	defer r.Body.Close()
	e := &event{}
	err = json.Unmarshal(body, e)
	if err != nil {
		log.Printf("error occured unmarshalling body: %v", err)
		json.NewEncoder(w).Encode("Error unmarshalling request")
		return
	}

	if e.Input != "" {
		e.Input = "hello: " + e.Input + ". How are you today?"
	}

	w.Write([]byte(e.Input))
}
```

You can see we have added several bits here, but the most important is the `HandlePost` function. This function will be called when a POST request is made to the `/event` route. This will read the request body and unmarshall it into a `event` struct and then modify the `Input` field of the struct to say hello to the name in the request we sent from the frontend.


## Build and Test

We can now build the frontend, update the static files, and run the application.

```
make debug
```

Navigate to `http://localhost:8080/` to see the form and response.


# Deploying to Kubernetes/Knative
Now that we have a working application, we can deploy it to Kubernetes very easily with the use of `ko`.

**Note** This also assumes you have Knative installed.

First lets create a manifest by creating a `deployment.yaml` file in a `config/` directory:
```
mkdir config
cd config && touch deployment.yaml
```

Populate this file with the following:
```
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: webapp
spec:
  template:
    spec:
      containers:
        - image: ko://github.com/JeffNeff/webapp/backend/cmd/webapp
```

If you want to apply this manifest to the `default` namespace, you can use the following command:

```
ko apply -f deployment.yaml
```

If you want to share a deployment manifest with a built image refrence, one can be created with

```
ko resolve -f deployment.yaml > share.yaml
```

# Deploying to a Container Hosting Service
If you would like to deploy via a container hosting service, Google Cloud Run for instance, you can use the following command to create the image:

```
ko resolve -f deployment.yaml > share.yaml
```

Then you can use the image located in the `share.yaml` file.

# Deploying Without ko
If you would like to deploy without using `ko`, you can create a `Dockerfile` in the `backend/` directory with the following contents:

```
FROM golang:1.17-buster AS builder
WORKDIR /project
COPY . ./
RUN cd /project/cmd/webapp && go build -o /project/bin/

FROM registry.access.redhat.com/ubi8/ubi-minimal
EXPOSE 8080
ENV KO_DATA_PATH /kodata
COPY --from=builder /project/cmd/webapp/kodata/ ${KO_DATA_PATH}/
COPY --from=builder /project/bin/webapp /webapp

ENTRYPOINT ["/webapp"]
```

Now one can build the image and deploy it to any service that supports hosting a container.
