import {getSiteUrl} from 'src/clients';
import {PARENT_ID_PARAM, SECTION_ID_PARAM} from 'src/constants';

export const isUrlHashValid = (urlHash: string, all: string[], some: string[]): boolean => {
    let valid = false;
    if (!urlHash || urlHash.length < 1) {
        return valid;
    }
    if (all.length > 0) {
        valid = all.every((a) => urlHash.includes(a));
    }
    if (some.length > 0) {
        valid = some.some((s) => urlHash.includes(s));
    }
    return valid;
};

export const isUrlEqualWithoutQueryParams = (url: string) => {
    const currentUrlWithoutQueryParams = window.location.href.split('?')[0];
    return currentUrlWithoutQueryParams === url || `${currentUrlWithoutQueryParams}/` === url;
};

export const isReferencedByUrlHash = (urlHash: string, id: string): boolean => {
    return urlHash === `#${id}`;
};

export const buildIdForUrlHashReference = (prefix: string, id: string): string => {
    return `${prefix}-${id}`;
};

export const buildToForCopy = (to: string): string => {
    return `${getSiteUrl()}${to}`;
};

export const buildTo = (
    fullUrl: string,
    id: string,
    query: string | undefined,
    url: string
) => {
    const isFullUrlProvided = fullUrl !== '';
    let to = isFullUrlProvided ? fullUrl : url;
    const isQueryProvided = query || query !== '';
    to = isQueryProvided ? `${to}?${query}` : to;
    return `${to}#${id}`;
};

export const buildQuery = (parentId: string, sectionId: string | undefined) => {
    let query = `${PARENT_ID_PARAM}=${parentId}`;
    if (sectionId) {
        query = `${query}&${SECTION_ID_PARAM}=${sectionId}`;
    }
    return query;
};