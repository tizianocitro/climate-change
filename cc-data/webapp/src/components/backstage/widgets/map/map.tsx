import React, {
    Dispatch,
    SetStateAction,
    useContext,
    useEffect,
    useState,
} from 'react';
import styled from 'styled-components';

import {AnchorLinkTitle, Header} from 'src/components/backstage/widgets/shared';
import {IsEcosystemRhsContext} from 'src/components/rhs/rhs_widgets';
import {FullUrlContext} from 'src/components/rhs/rhs';
import {buildQuery, useScrollIntoView, useUrlHash} from 'src/hooks';
import {formatName} from 'src/helpers';
import {MapData, Point, defaultMapData} from 'src/types/map';
import {Spacer} from 'src/components/backstage/grid';

import WorldMap from './world_map';
import PointSelect from './point_select';

type Props = {
    data: MapData;
    name: string;
    parentId: string;
    sectionId: string;
    point: string;
    setPoint: Dispatch<SetStateAction<string>>;
};

const Map = ({
    data,
    name = '',
    parentId,
    sectionId,
    point,
    setPoint,
}: Props) => {
    const isEcosystemRhs = useContext(IsEcosystemRhsContext);
    const fullUrl = useContext(FullUrlContext);
    const urlHash = useUrlHash();

    const {items, points, range, colorRange, worldEnv, seaEnv} = data;

    const [selectedPoint, setSelectedPoint] = useState<Point>(points.defaultPoint || defaultMapData.points.defaultPoint);
    useEffect(() => {
        setPoint(selectedPoint.value);
    }, [selectedPoint]);
    useEffect(() => {
        setSelectedPoint({label: point, value: point});
    }, [point]);

    const id = `${formatName(name)}-${sectionId}-${parentId}-widget`;
    const ecosystemQuery = isEcosystemRhs ? '' : buildQuery(parentId, sectionId);

    useScrollIntoView(urlHash);

    return (
        <Container
            id={id}
            data-testid={id}
        >
            <Header>
                <AnchorLinkTitle
                    fullUrl={fullUrl}
                    id={id}
                    query={ecosystemQuery}
                    text={name}
                    title={name}
                />
                <Spacer/>
                <PointSelect
                    data={points}
                    selectedPoint={selectedPoint}
                    setSelectedPoint={setSelectedPoint}
                    setPoint={setPoint}
                />
            </Header>
            <WorldMap
                countries={items}
                range={range}
                colorRange={colorRange}
                worldEnv={worldEnv}
                seaEnv={seaEnv}
                selectedPoint={selectedPoint}
                parentId={parentId}
                sectionId={sectionId}
            />
        </Container>
    );
};

const Container = styled.div`
    width: 100%;
    display: flex;
    flex-direction: column;
    margin-top: 24px;
`;

export default Map;