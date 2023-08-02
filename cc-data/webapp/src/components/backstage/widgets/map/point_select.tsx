import React, {
    Dispatch,
    FC,
    SetStateAction,
    useContext,
    useEffect,
} from 'react';
import {Button, Select} from 'antd';
import {ArrowLeftOutlined, ArrowRightOutlined} from '@ant-design/icons';
import {useIntl} from 'react-intl';
import styled from 'styled-components';

import {IsRhsContext} from 'src/components/backstage/sections_widgets/sections_widgets_container';
import {Point, PointData} from 'src/types/map';
import {useUrlHash} from 'src/hooks';

type Props = {
    data: PointData;
    selectedPoint: Point;
    setSelectedPoint: Dispatch<SetStateAction<Point>>;
    setPoint: Dispatch<SetStateAction<string>>;
};

export const getPointFromUrlHash = (urlHash: string): string => {
    if (urlHash.includes('mapel-')) {
        return urlHash.substring(7).split('-')[0];
    }
    if (urlHash.includes('sea-')) {
        // It's sea env
        return urlHash.substring(5).split('-')[0];
    }
    return urlHash;
};

export const isInPoints = (data: PointData, point: string): boolean => data.points.some((p) => p.value === point);

const PointSelect: FC<Props> = ({
    data,
    selectedPoint,
    setSelectedPoint,
    setPoint,
}) => {
    const isRhs = useContext(IsRhsContext);
    const {formatMessage} = useIntl();
    const urlHash = useUrlHash();

    useEffect(() => {
        if (!urlHash || urlHash.length < 1) {
            return;
        }
        const point = getPointFromUrlHash(urlHash);
        if (isInPoints(data, point)) {
            setSelectedPoint({value: point, label: point});
        }
    }, [urlHash]);

    useEffect(() => {
        setPoint(selectedPoint.value);
    }, [selectedPoint]);

    const filterSort = (optionA: Point, optionB: Point): number => {
        return (optionA?.label ?? '').toLowerCase().localeCompare((optionB?.label ?? '').toLowerCase());
    };

    const filterOption = (input: string, option: Point | undefined): boolean => {
        return (option?.label ?? '').includes(input);
    };

    const setNextPointAsSelected = () => {
        if (!selectedPoint || !selectedPoint.value) {
            return;
        }
        const value = `${parseInt(selectedPoint.value, 10) + 1}`;
        setSelectedPoint({value, label: value});
    };

    const setPreviousPointAsSelected = () => {
        if (!selectedPoint || !selectedPoint.value) {
            return;
        }
        const value = `${parseInt(selectedPoint.value, 10) - 1}`;
        setSelectedPoint({value, label: value});
    };

    const isButtonDisabled = (isPrev: boolean): boolean => {
        if (!selectedPoint || !selectedPoint.value || points.length < 1) {
            return true;
        }
        if (isPrev) {
            return selectedPoint.value === points[0].value;
        }
        return selectedPoint.value === points[points.length - 1].value;
    };

    const {points} = data;
    const width = isRhs ? 125 : 200;
    return (
        <Container>
            <Button
                style={{marginRight: '1%'}}
                key='prev'
                onClick={setPreviousPointAsSelected}
                disabled={isButtonDisabled(true)}
                icon={<ArrowLeftOutlined/>}
            />
            <Select
                value={selectedPoint.value}
                showSearch={true}
                style={{width}}
                placeholder={formatMessage({defaultMessage: 'Search or select a year'})}
                optionFilterProp='children'
                filterOption={filterOption}
                filterSort={filterSort}
                options={points}
                onChange={(value) => setSelectedPoint({value, label: value})}
            />
            <Button
                style={{marginLeft: '1%'}}
                key='next'
                onClick={setNextPointAsSelected}
                disabled={isButtonDisabled(false)}
                icon={<ArrowRightOutlined/>}
            />
        </Container>
    );
};

const Container = styled.div`
    display: flex;
`;

export default PointSelect;