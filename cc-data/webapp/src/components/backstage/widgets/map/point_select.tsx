import React, {
    Dispatch,
    FC,
    SetStateAction,
    useContext,
    useEffect,
} from 'react';
import {Select} from 'antd';
import {useIntl} from 'react-intl';

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

    const {points} = data;
    const width = isRhs ? 125 : 200;
    return (
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
    );
};

export default PointSelect;