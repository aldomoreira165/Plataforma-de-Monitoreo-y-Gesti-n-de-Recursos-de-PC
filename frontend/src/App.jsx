import Navegacion from "./componentes/Navegacion"
import RamCPU from "./paginas/RamCPU"
import RamCPUH from "./paginas/RamCPUH" 
import ArbolProcesos from "./paginas/ArbolProcesos";
import SimulacionEstadosProcesos from "./paginas/SimulacionEstadosProcesos";
import { BrowserRouter as Router, Routes, Route, NavLink } from "react-router-dom";

function App() {
  return (
    <Router>
      <Navegacion />
      <Routes>
        <Route exact path="/" Component={ RamCPU }></Route>
        <Route exact path="/RAMCPUH" Component={ RamCPUH }></Route>
        <Route exact path="/ArbolProcesos" Component={ ArbolProcesos }></Route>
        <Route exact path="/SimulacionEstadosProcesos" Component={ SimulacionEstadosProcesos }></Route>
      </Routes>
    </Router>
  )
}

export default App
