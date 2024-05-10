import { Line } from 'react-chartjs-2';
import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend,
    Filler,
} from 'chart.js';

export default function LineChart(props) {

    ChartJS.register(
        CategoryScale,
        LinearScale,
        PointElement,
        LineElement,
        Title,
        Tooltip,
        Legend,
        Filler
    );

    var porcentajes_uso = props.y;
    var tiempos = props.x;

    var midata = {
        labels: tiempos,
        datasets: [
            {
                label: 'Porcentaje Usado',
                data: porcentajes_uso,
                borderColor: 'rgb(255, 99, 132)',
                backgroundColor: 'rgba(255, 99, 132, 0.5)',
                pointBorderColor: 'rgba(255, 99, 132)',
                pointBackgroundColor: 'rgba(255, 99, 132)',
            }, {
                label: 'Porcentaje Libre',
                data: props.libre,
                borderColor: 'rgb(54, 162, 235)',
                backgroundColor: 'rgba(54, 162, 235, 0.5)',
                pointBorderColor: 'rgba(54, 162, 235)',
                pointBackgroundColor: 'rgba(54, 162, 235)',
            },
        ],
    };

    var misoptions = {
        scales : {
            y : {
                min : 0
            }
        }
    };

    return <Line data={midata} options={misoptions}/>
}