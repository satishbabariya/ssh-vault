import { useAuth0 } from "@auth0/auth0-react";
import { Button, Center } from "@mantine/core";
import React from "react";
import { Loader } from "./Loader";

export function Auth() {
  const { isLoading, isAuthenticated, error, user, loginWithRedirect } =
    useAuth0();

  if (isLoading) {
    return <Loader />;
  }

  return (
    <Center
      style={{
        height: "100vh",
      }}
    >
      <Button onClick={loginWithRedirect}>Log in</Button>
    </Center>
  );
}
