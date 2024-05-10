import React from 'react';
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from "chart.js";
import { Pie } from 'react-chartjs-2';

ChartJS.register(ArcElement, Tooltip, Legend);

export default function Pies(props) {

    var options = {
        responsive : true,
        maintainAspectRatio: false,
    };
    
    var data = {
        labels: ['Libre', 'Usada'],
        datasets: [
            {
                data: [props.libre, props.usada],
                backgroundColor: [
                    'rgba(54, 162, 235, 0.2)',
                    'rgba(255, 99, 132, 0.2)',
                ],
                borderColor: [
                    'rgba(54, 162, 235)',
                    'rgba(255, 99, 132, 1)',
                ],
                borderWidth: 3,
            },
        ],
    };

    return (
        <>
            <Pie data={data} options={options} />
        </>
    ); 
}