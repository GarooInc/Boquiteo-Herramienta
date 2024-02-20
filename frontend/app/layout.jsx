import "@/styles/globals.css";

export const metadata = {
    title: 'Boquiteo',
}

const RootLayout = ({children}) => {
  return (
    <html lang="es">
        <head>
            <title>{metadata.title}</title>
            <link rel="icon" type="image/png+xml" href="/assets/images/biglogo03.png" />
            <meta name="viewport" content="width=device-width, initial-scale=1.0" />
            <meta property="og:title" content={metadata.title} />
        </head>
        <body>
            <main className='app'>
                {children}
            </main>
        </body>
    </html>
  )
}

export default RootLayout