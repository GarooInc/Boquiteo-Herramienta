import React from 'react'

const principalPage = () => {
    const menu = [
        "Call Center",
        "Cocina",
        "Entrega",
    ]

  return (
    <div className='w-full isolate h-screen'>
        <div className='flex h-full'>
            <div className='w-1/3 h-full'>
                <img src='/assets/logoboquiteo.png' alt='logo' className='w-1/2 m-auto' />
                <div className='flex flex-col gap-6'>
                {
                    menu.map((item, index) => {
                        return (
                            <div key={index} className='w-full px-8'>
                                <a href={item.toLowerCase()} className='block text-center text-black hover:bg-orange hover:text-white p-2 font-bold border-2 border-orange rounded-md'>{item}</a>
                            </div>
                        )
                    })
                }
                </div>
            </div>
            <div className='w-full bg-light-gray h-full'>
            </div>
        </div>
    </div>
  )
}

export default principalPage