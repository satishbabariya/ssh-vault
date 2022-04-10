import React, { useEffect } from "react";
import { useAuth0 } from "@auth0/auth0-react";
import { Credentials } from "./components/Credentials";
import { Loader } from "./components/Loader";
import { Auth } from "./components/Auth";
import { showNotification } from "@mantine/notifications";

export function App() {
  const { isLoading, error, isAuthenticated } = useAuth0();

  useEffect(() => {
    if (error) {
      showNotification({
        message: error.message,
        color: "red",
      });
    }
  }, [error]);

  if (isLoading) {
    return <Loader />;
  }

  if (isAuthenticated) {
    return <Credentials />;
  } else {
    return <Auth />;
  }
}
