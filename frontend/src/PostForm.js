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
