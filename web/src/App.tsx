import { useAuth0 } from "@auth0/auth0-react";
import { Button, Center, Stack, Text } from "@mantine/core";
import React from "react";
import { Loader } from "./components/Loader";

export default function App() {
  const { isLoading, isAuthenticated, error, user, loginWithRedirect, logout } =
    useAuth0();

  console.log(user, error);
  if (isLoading) {
    return <Loader />;
  }
  if (isAuthenticated) {
    return (
      <Center
        style={{
          height: "100vh",
        }}
      >
        <Stack>
          <Text>Hello {user!.name}</Text>
          <Center>
            <Button
              onClick={() => logout({ returnTo: window.location.origin })}
            >
              Log out
            </Button>
          </Center>
        </Stack>
      </Center>
    );
  } else {
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
}
