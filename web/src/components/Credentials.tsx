import { useAuth0 } from "@auth0/auth0-react";
import { Button, Center, Stack, Table, Text } from "@mantine/core";
import React from "react";
import { Auth0Config } from "../config";
import { useApi } from "../hooks/use-api";
import { Loader } from "./Loader";

export function Credentials() {
  const {
    loading,
    error,
    data: remotes = [],
  } = useApi(`/api/credentials`, {
    audience: Auth0Config.audience,
  });

  if (loading) {
    return <Loader />;
  }

  if (error) {
    return (
      <Center
        style={{
          height: "100vh",
        }}
      >
        <Text>{error.message}</Text>
      </Center>
    );
  }

  return (
    <Center>
      <Table>
        <thead>
          <tr>
            <th>Server</th>
            <th>PORT</th>
          </tr>
        </thead>
        <tbody>
          {remotes.map(({ host, port }: any) => (
            <tr key={`${host}:${port}`}>
              <td>{host}</td>
              <td>{port}</td>
            </tr>
          ))}
        </tbody>
      </Table>
    </Center>
  );
}
