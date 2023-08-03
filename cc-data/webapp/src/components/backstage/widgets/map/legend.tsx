import React, {FC} from 'react';
import styled from 'styled-components';

export type LegendData = {
    minValue: number;
    maxValue: number;
    minColor: string;
    maxColor: string;
};

type Props = {
    data: LegendData;
};

const Legend: FC<Props> = ({data}) => {
    const {minValue, maxValue, minColor, maxColor} = data;

    // Calculate the percentage based on the input values
    const percentage = (value: number) => ((value - minValue) / (maxValue - minValue)) * 100;

    return (
        <LegendContainer>
            <GradientLine
                style={{
                    background: `linear-gradient(to right, ${minColor} ${percentage(minValue)}%, ${maxColor} ${percentage(maxValue)}%)`,
                }}
            />
        </LegendContainer>
    );
};

const LegendContainer = styled.div`
    width: 60%;
    height: 20px;
    background-color: #f0f0f0;
    border-radius: 5px;
    overflow: hidden;
`;

const GradientLine = styled.div`
    width: 100%;
    height: 100%;
    border-radius: 5px;
`;

export default Legend;
