import { React, useState, useEffect } from 'react'
import PieChart from '../componentes/PieChart'
import '../assets/css/main.css'

export default function RamCPU() {

  const [ramLibre, setRamLIbre] = useState(0)
  const [ramUsada, setRamUsada] = useState(0)
  const [cpuLibre, setCpuLIbre] = useState(0)
  const [cpuUsada, setCpuUsada] = useState(0)

  const handleDataRam = async () => {
    const res = await fetch('/api/ram')
    console.log(res)
    const data = await res.json()
    setRamLIbre(data.libre)
    setRamUsada(data.usada)

    // Enviar datos a la base de datos
    insertarUsoRam(data.usada);
  }

  const handleDataCpu = async () => {
    const res = await fetch('/api/cpu')
    console.log(res)
    const data = await res.json()
    setCpuLIbre(100 - data.usada)
    setCpuUsada(data.usada)

    // Enviar datos a la base de datos
    insertarUsoCpu(data.usada);
  }

  const insertarUsoRam = async (usada) => {
    const tiempo = new Date().toISOString();
    const requestOptions = {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ porcentaje: usada, tiempo: tiempo }),
    };
    await fetch('/api/insertar_uso_ram', requestOptions);
  };

  const insertarUsoCpu = async (usada) => {
    const tiempo = new Date().toISOString();
    const requestOptions = {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ porcentaje: usada, tiempo: tiempo }),
    };
    await fetch('/api/insertar_uso_cpu', requestOptions);
  };

  useEffect(() => {
    // obteniendo y cargando datos
    const interval = setInterval(() => {

      handleDataRam();
      handleDataCpu();
    }, 1000);

    return () => clearInterval(interval);
  }, []);

  return (
    <>
      <div className='main-container'>
        <h1 className='main-container__title'>Monitoreo en Tiempo Real</h1>
        <div className='chart-container'>
          <div className='chart-container__titles'>
            <h3>RAM</h3>
            <h3>CPU</h3>
          </div>
          <div className='chart-container__elements'>
            <div className='ram-container'>
              <PieChart libre={ramLibre} usada={ramUsada} />
            </div>
            <div className='cpu-container'>
              <PieChart libre={cpuLibre} usada={cpuUsada} />
            </div>
          </div>
        </div>
      </div>
    </>
  )
}
