import React from 'react'
import { useNavigate } from 'react-router-dom'

export const BackButton = ({url}) => {
    const navigate = useNavigate()
  return (
    <button className='back-button' onClick={()=>navigate(url)}>back</button>
  )
}
