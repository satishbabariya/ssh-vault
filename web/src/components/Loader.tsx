import { Center, Loader as MantieLoader } from "@mantine/core";
import React from "react";

export function Loader() {
  return (
    <Center
      style={{
        height: "100vh",
      }}
    >
      <MantieLoader />
    </Center>
  );
}
