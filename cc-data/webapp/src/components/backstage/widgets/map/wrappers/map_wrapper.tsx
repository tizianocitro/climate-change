import React, {useContext} from 'react';
import {useLocation, useRouteMatch} from 'react-router-dom';
import qs from 'qs';

import {SectionContext} from 'src/components/rhs/rhs';
import Map from 'src/components/backstage/widgets/map/map';
import {MapData} from 'src/types/map';

type Props = {
    name?: string;
    url?: string;
};

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

    // const data = useListData(formatUrlWithId(url, sectionIdForUrl));
    const data: MapData = {
        items: [
            {id: 'usa', iso3: 'USA', country: 'United States', value: 100},
            {id: 'canada', iso3: 'CAN', country: 'Canada', value: 200},
            {id: 'italy', iso3: 'ITA', country: 'Italy', value: 150},
            {id: 'brazil', iso3: 'BRA', country: 'Brazil', value: 180},
        ],
        points: {
            defaultPoint: {
                value: '2022',
                label: '2022',
            },
            points: [
                {
                    value: '2020',
                    label: '2020',
                },
                {
                    value: '2021',
                    label: '2021',
                },
                {
                    value: '2022',
                    label: '2022',
                },
            ],
        },
    };

    return (
        <Map
            data={data}
            name={name}
            sectionId={sectionIdForUrl}
            parentId={parentId}
        />
    );
};

export default MapWrapper;