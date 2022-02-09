import type { AppProps } from 'next/app'
import { ChakraProvider, extendTheme, theme as defaultTheme } from '@chakra-ui/react'
import Navbar from '@components/shared/Navbar'
import '@fontsource/ubuntu'
import '@fontsource/open-sans'
import { Auth0Provider } from '@auth0/auth0-react'

const auth0ClientID = process.env.NEXT_PUBLIC_AUTH0_CLIENT_ID
const auth0Domain = process.env.NEXT_PUBLIC_AUTH0_DOMAIN
const appDomain = process.env.NEXT_PUBLIC_APP_DOMAIN


if(!auth0ClientID || !auth0Domain || !appDomain){
  throw new Error("error: AUTH0 configuration is not correct")
}

const theme = extendTheme({
  colors: {
    primary: defaultTheme.colors.gray,
    secondary: defaultTheme.colors.gray,
    text: defaultTheme.colors.gray
  },
  fonts: {
    heading: 'Ubuntu',
    text: 'Open Sans'
  }
})

const App = ({ Component, pageProps }: AppProps) => {
  return (
    <ChakraProvider theme={theme} resetCSS>
      <Auth0Provider
      domain={auth0Domain}
      clientId={auth0ClientID}
      redirectUri={appDomain}>
      <Navbar>
        <Component {...pageProps} />
      </Navbar>
      </Auth0Provider>
    </ChakraProvider>
  )
}

export default App
