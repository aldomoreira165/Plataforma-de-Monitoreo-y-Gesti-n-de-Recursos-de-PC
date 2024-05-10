import { React, useState, useEffect } from 'react';
import LinesChart from '../componentes/LineChart';
import '../assets/css/main.css'

export default function RamCPUH() {
  const [usoRamData, setUsoRamData] = useState([]);
  const [usoCpuData, setUsoCpuData] = useState([]);

  useEffect(() => {
    // obteniendo los datos historicos de la ram
    fetch('/api/uso_historico_ram')
      .then(response => response.json())
      .then(data => {
        console.log(data);
        setUsoRamData(data);
      })
      .catch(error => console.log(error));

    // obteniendo los datos historicos de la cpu
    fetch('/api/uso_historico_cpu')
      .then(response => response.json())
      .then(data => {
        console.log(data);
        setUsoCpuData(data);
      })
      .catch(error => console.log(error));
  }, []);

  // Procesamiento de datos de RAM
  const eje_x_ram = usoRamData.map(item => item.tiempo);
  const eje_y_ram = usoRamData.map(item => item.porcentaje);
  const porcentajes_libres_ram = eje_y_ram.map(porcentaje => 100 - porcentaje);

  // Procesamiento de datos de CPU
  const eje_x_cpu = usoCpuData.map(item => item.tiempo);
  const eje_y_cpu = usoCpuData.map(item => item.porcentaje);
  const porcentajes_libres_cpu = eje_y_cpu.map(porcentaje => 100 - porcentaje);

  return (
    <>
      <div className='main-container'>
        <h1 className='main-container__title'>Monitoreo Histórico</h1>
        <div className='chart-container'>
          <div className='chart-container__titles'>
            <h3>RAM</h3>
            <h3>CPU</h3>
          </div>
          <div className='chart-container__elements'>
            <div className='ram-container'>
              <LinesChart x={eje_x_ram} y={eje_y_ram} libre={porcentajes_libres_ram} />
            </div>
            <div className='cpu-container'>
              <LinesChart x={eje_x_cpu} y={eje_y_cpu} libre={porcentajes_libres_cpu} />
            </div>
          </div>
        </div>
      </div>
    </>
  )
}
