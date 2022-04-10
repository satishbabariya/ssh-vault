import { useAuth0 } from "@auth0/auth0-react";
import { Button, Center, Stack, Text } from "@mantine/core";
import React from "react";
import { Loader } from "./Loader";

export function Credentials() {
  const { isLoading, error, logout, user } = useAuth0();

  if (isLoading) {
    return <Loader />;
  }

  return (
    <Center
      style={{
        height: "100vh",
      }}
    >
      <Stack>
        <Text>Hello {user!.name}</Text>
        <Center>
          <Button onClick={() => logout({ returnTo: window.location.origin })}>
            Log out
          </Button>
        </Center>
      </Stack>
    </Center>
  );
}
