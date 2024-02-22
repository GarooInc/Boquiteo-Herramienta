import React from 'react'
import { useRouter, usePathname } from "next/navigation"



const NavBar = () => {
    const router = useRouter()
    const pathname = usePathname()


    const menu = [
        { name: "Call Center", link: "/" },
        { name: "Cocina", link: "/cocina" },
        { name: "Entrega", link: "/entrega" },
    ]

    const navigateTo = (link) => {
        console.log(link)
        router.push(link)
    }

    // Función para determinar si el enlace es la ruta activa
    const isActive = (link) => {
        return link === pathname
    }

  return (
    <div className='lg:w-1/3 lg:h-full w-full'>
        <img src='/assets/logoboquiteo.png' alt='logo' className='w-1/2 m-auto' />
        <div className='flex flex-col gap-6'>
            {
                menu.map((item, index) => {
                    const buttonClass = isActive(item.link)
                        ? 'w-full block text-center text-white bg-orange p-2 font-bold border-2 border-orange rounded-md'
                        : 'w-full block text-center text-black hover:bg-orange hover:text-white p-2 font-bold border-2 border-orange rounded-md';
                    return (
                        <div key={index} className='w-full px-8'>
                            <button className={buttonClass} onClick={() => navigateTo(item.link)}>{item.name}</button>
                        </div>
                    )
                })
            }
        </div>
    </div>
  )
}

export default NavBar
