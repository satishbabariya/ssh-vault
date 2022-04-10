import React from "react";
import ReactDOM from "react-dom";
import App from "./App";

import { Auth0Provider, AppState } from "@auth0/auth0-react";
import { MantineProvider } from "@mantine/core";
import { Auth0Config } from "./config";

const onRedirectCallback = (appState?: AppState) => {
  window.history.replaceState({}, document.title, window.location.pathname);
  // history.replace((appState && appState.returnTo) || window.location.pathname);
};

ReactDOM.render(
  <React.StrictMode>
    <Auth0Provider
      domain={Auth0Config.domain}
      clientId={Auth0Config.clientId}
      redirectUri={window.location.origin}
      onRedirectCallback={onRedirectCallback}
    >
      <MantineProvider
        theme={{
          primaryColor: "dark",
          defaultRadius: "xs",
          // fontFamily: "Greycliff CF, sans-serif",
        }}
      >
        <App />
      </MantineProvider>
    </Auth0Provider>
  </React.StrictMode>,
  document.getElementById("root")
);
