"use client";
import React, {useState, useEffect} from 'react'
import NavBar from '@/components/NavBar/NavBar';


const Page = () => {
    const [data , setData] = useState([])
    const [searchTerm, setSearchTerm] = useState('')

    const filteredData = data.filter(item =>
        item.customer.toLowerCase().includes(searchTerm.toLowerCase()) ||
        item.order_number.toString().includes(searchTerm)
    )
       

    const fetchData = async () => {
        const response = await fetch("https://boquiteo-garoo.koyeb.app/api/orders")
        const data = await response.json()
        if (data.status === 200) {
            setData(data.data)
        }
        else {
            console.log("Error")
        }
        console.log(data)
    }

    useEffect(() => {
        fetchData()
    }
    , [])

  return (
    <div className='w-full isolate h-screen'>
        <div className='flex h-full lg:flex-row flex-col'>
            <NavBar />
            <div className='w-full bg-light-gray h-full flex flex-col justify-start items-center'>
                <input
                    type='text'
                    placeholder='Buscar orden'
                    className='w-2/3 m-4 p-2 rounded-md'
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                />
                <div className='flex flex-wrap justify-center gap-4 w-full overflow-y-auto'>
                    {
                        filteredData.map((item) => {
                            return (
                                <div key={item._id} className='lg:w-1/4 lg:m-0 mx-4 w-full bg-white p-4 rounded-md border-2 border-dark-gray flex flex-col items-start justify-start gap-2'>
                                    <div className='flex flex-col'>
                                        <h1 className='text-xl font-bold'>Order #{item.order_number}</h1>
                                        <p className='text-l text-dark-gray uppercase'>{item.customer}</p>
                                    </div>
                                    <p className='text-lg'>{item.description}</p>
                                    <p className='text-lg'>{item.price}</p>
                                    {
                                        item.line_items.map((line_item, index) => {
                                            return (
                                                <div key={index} className='flex flex-col border-b-2 border-dark-gray w-full p-2'>
                                                    <p className='text-lg'>{line_item.name}</p>
                                                    <p className='text-lg text-gray'>Cocina: {line_item.vendor}</p>
                                                    <div className='flex flex-row justify-between'>
                                                        <p className='text-lg'>Qty: {line_item.quantity}</p>
                                                        <p className='text-lg font-bold'>Q{line_item.price}</p>
                                                    </div>
                                                    <p className='text-lg'>{line_item.status}</p>
                                                </div>
                                            )

                                        }
                                        )
                                        
                                    }
                                </div>
                            )
                        }
                    )      
                    }
                </div>
            </div>
        </div>
    </div>
  )
}

export default Page