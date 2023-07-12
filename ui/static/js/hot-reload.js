function hotReload() {
  let socket = new WebSocket("ws://localhost:4000/reload");

  socket.onopen = () => {
    console.debug("Websocket connection successful. Hot reloading active");
  };

  socket.onclose = (event) => {
    console.debug("Hot reload connection closed. Waiting for connection...");
    let counter = 0;
    const timer = setInterval(function () {
      if (counter >= 20) {
        clearInterval(timer);
        console.debug(
          "Hot reload connection could not be made. Reload page to establish again"
        );
      }
      fetch("http://localhost:4000/reload-ready")
        .then((res) => {
          if (res.status === 200) {
            location.replace(location.href);
          }
        })
        .catch((err) => {});
      counter++;
    }, 500);
  };

  socket.onerror = (error) => {
    console.debug("Socket error: ", error);
    socket.close();
  };
}

hotReload();
