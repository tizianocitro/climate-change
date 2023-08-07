import {useEffect} from 'react';
import {useSelector} from 'react-redux';
import {useLocation} from 'react-router-dom';
import {getCurrentChannelId, getCurrentUser, getCurrentUserId} from 'mattermost-webapp/packages/mattermost-redux/src/selectors/entities/common';
import {getCurrentTeamId} from 'mattermost-webapp/packages/mattermost-redux/src/selectors/entities/teams';

import {saveUrlHashTelemetry} from 'src/clients';
import {channelNameSelector, teamNameSelector} from 'src/selectors';

import {isUrlHashValid} from './url';

export const useUrlHash = (): string => {
    const {hash: urlHash} = useLocation();
    let renderHash = localStorage.getItem('previousHash') || '';
    renderHash = urlHash && urlHash !== '' ? urlHash : renderHash;
    return renderHash;
};

export const useCleanUrlHash = () => {
    useEffect(() => {
        const hash = localStorage.getItem('previousHash');
        if (!hash) {
            localStorage.setItem('previousHash', '');
            return;
        }
        const element = document.querySelector(hash);
        if (!element) {
            localStorage.setItem('previousHash', '');
        }
    });
};

export const useCleanUrlHashOnChannelChange = () => {
    const channelId = useSelector(getCurrentChannelId);
    useEffect(() => {
        localStorage.setItem('previousHash', '');
    }, [channelId]);
};

export const useSaveUrlHashTelemetry = (hash: string) => {
    const channelId = useSelector(getCurrentChannelId);
    const channel = useSelector(channelNameSelector(channelId));
    const teamId = useSelector(getCurrentTeamId);
    const team = useSelector(teamNameSelector(teamId));
    const userId = useSelector(getCurrentUserId);
    const user = useSelector(getCurrentUser);

    useEffect(() => {
        if (!isUrlHashValid(hash, [], [])) {
            return;
        }
        async function saveUrlHashTelemetryAsync() {
            saveUrlHashTelemetry({
                channelId,
                channelName: channel?.display_name || '',
                teamId,
                teamName: team?.display_name || '',
                userId,
                username: user?.username || '',
                urlHash: hash,
            });
        }
        saveUrlHashTelemetryAsync();
    }, [hash]);
};

type ScrollIntoViewPositions = {
    block?: ScrollLogicalPosition;
    inline?: ScrollLogicalPosition;
};

export const useScrollIntoView = (hash: string, positions?: ScrollIntoViewPositions) => {
    useCleanUrlHash();

    // When first loading the page, the element with the ID corresponding to the URL
    // hash is not mounted, so the browser fails to automatically scroll to such section.
    // To fix this, we need to manually scroll to the component
    useEffect(() => {
        const options = buildOptions(positions);
        const previousHash = localStorage.getItem('previousHash');
        if (hash !== '' || previousHash) {
            setTimeout(() => {
                let urlHash = hash;
                if (urlHash === '' && previousHash) {
                    urlHash = previousHash;
                }
                document.querySelector(urlHash)?.scrollIntoView(options);
                localStorage.setItem('previousHash', urlHash);
                window.location.hash = '';
            }, 300);
        }
    }, [hash]);

    useCleanUrlHashOnChannelChange();
    useSaveUrlHashTelemetry(hash);
};

// Doc: https://developer.mozilla.org/en-US/docs/Web/API/Element/scrollIntoView
const buildOptions = (positions: ScrollIntoViewPositions | undefined): ScrollIntoViewOptions => {
    let options: ScrollIntoViewOptions = {
        behavior: 'smooth',
        block: 'center',
        inline: 'nearest',
    };
    if (positions) {
        const {block, inline} = positions;
        options = {...options, block, inline};
    }
    return options;
};

// export const useScrollIntoViewWithCustomTime = (hash: string, time: number) => {
//     useEffect(() => {
//         if (hash !== '') {
//             setTimeout(() => {
//                 document.querySelector(hash)?.scrollIntoView({behavior: 'smooth'});
//             }, time);
//         }
//     }, [hash]);
// };