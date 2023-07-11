import React, {useContext, useEffect, useState} from 'react';
import {useLocation, useRouteMatch} from 'react-router-dom';
import qs from 'qs';

import {SectionContext} from 'src/components/rhs/rhs';
import Map from 'src/components/backstage/widgets/map/map';
import {useMapData, useUrlHash} from 'src/hooks';
import {formatUrlWithId} from 'src/helpers';
import {defaultMapData} from 'src/types/map';

import {getPointFromUrlHash, isInPoints} from 'src/components/backstage/widgets/map/point_select';

type Props = {
    name?: string;
    url?: string;
};

const replaceQueryParams = (url: string, year?: string): string => url.replace(':year', year || '2022');

const MapWrapper = ({
    name = '',
    url = '',
}: Props) => {
    const sectionContextOptions = useContext(SectionContext);
    const {params: {sectionId}} = useRouteMatch<{sectionId: string}>();
    const location = useLocation();
    const queryParams = qs.parse(location.search, {ignoreQueryPrefix: true});
    const parentIdParam = queryParams.parentId as string;

    const areSectionContextOptionsProvided = sectionContextOptions.parentId !== '' && sectionContextOptions.sectionId !== '';
    const parentId = areSectionContextOptionsProvided ? sectionContextOptions.parentId : parentIdParam;
    const sectionIdForUrl = areSectionContextOptionsProvided ? sectionContextOptions.sectionId : sectionId;

    const urlHash = useUrlHash();
    const [point, setPoint] = useState<string>(defaultMapData.points.defaultPoint.value);

    const data = useMapData(replaceQueryParams(formatUrlWithId(url, sectionIdForUrl), point));

    useEffect(() => {
        if (!urlHash || urlHash.length < 1) {
            return;
        }
        const pointFromUrlHash = getPointFromUrlHash(urlHash);
        if (isInPoints(data.points || defaultMapData.points, pointFromUrlHash)) {
            setPoint(pointFromUrlHash);
        }
    });

    // const data: MapData = {
    //     items: [
    //         {id: 'usa', iso3: 'USA', country: 'United States', value: 100},
    //         {id: 'canada', iso3: 'CAN', country: 'Canada', value: 200},
    //         {id: 'italy', iso3: 'ITA', country: 'Italy', value: 150},
    //         {id: 'brazil', iso3: 'BRA', country: 'Brazil', value: 180},
    //     ],
    //     points: {
    //         defaultPoint: {
    //             value: '2022',
    //             label: '2022',
    //         },
    //         points: [
    //             {
    //                 value: '2020',
    //                 label: '2020',
    //             },
    //             {
    //                 value: '2021',
    //                 label: '2021',
    //             },
    //             {
    //                 value: '2022',
    //                 label: '2022',
    //             },
    //         ],
    //     },
    // };

    return (
        <>
            {data &&
                <Map
                    data={data}
                    name={name}
                    sectionId={sectionIdForUrl}
                    parentId={parentId}
                    point={point}
                    setPoint={setPoint}
                />}
        </>
    );
};

export default MapWrapper;