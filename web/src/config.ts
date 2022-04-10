export const Auth0Config = {
  domain: import.meta.env.VITE_APP_DOMAIN as string,
  clientId: import.meta.env.VITE_APP_CLIENT_ID as string,
  audience: import.meta.env.VITE_APP_AUDIENCE as string,
};
