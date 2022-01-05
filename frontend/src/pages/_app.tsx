import type { AppProps } from 'next/app'
import { ChakraProvider, extendTheme, theme as defaultTheme } from '@chakra-ui/react'
import Navbar from '@components/shared/Navbar'
import '@fontsource/ubuntu'
import '@fontsource/open-sans'

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
      <Navbar>
        <Component {...pageProps} />
      </Navbar>
    </ChakraProvider>
  )
}

export default App
