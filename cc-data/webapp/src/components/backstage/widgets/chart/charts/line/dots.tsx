import React, {
    FC,
    useContext,
    useEffect,
    useState,
} from 'react';
import {useIntl} from 'react-intl';
import {useRouteMatch} from 'react-router-dom';
import styled from 'styled-components';

import {formatUrlAsMarkdown} from 'src/components/backstage/header/controls';
import {useToaster} from 'src/components/backstage/toast_banner';
import {FullUrlContext} from 'src/components/rhs/rhs';
import {IsEcosystemRhsContext} from 'src/components/rhs/rhs_widgets';
import {buildQuery, buildTo, buildToForCopy} from 'src/hooks';
import {copyToClipboard} from 'src/utils';

export const DOT_PREFIX = 'dot-';

const dotStringify = (x: string, y: string, value: string | number): [string, string, string] => {
    const cxString = `${x}`.replace('.', 'dot');
    const cyString = `${y}`.replace('.', 'dot');
    const valueString = `${value}`.replace('.', 'dot');
    return [cxString, cyString, valueString];
};

export const valueStringify = (value: string | number): string => {
    return `${value}`.replace('.', 'dot');
};

export const idStringify = (id: string): string => {
    return `${id}`.replaceAll('-', 'hpn');
};

type Props = any;

export const Dot: FC<Props> = (props) => {
    const {cx, cy, fill, payload, value, originalColor, selectedDot, sectionId} = props;
    const [color, setColor] = useState(fill);

    useEffect(() => {
        // if (selectedDot.x === cx && selectedDot.y === cy) {
        if (selectedDot.label === payload.label && valueStringify(selectedDot.value) === valueStringify(value)) {
            setColor('#F4B400');
        } else {
            setColor(originalColor);
        }
    }, [selectedDot]);

    return (
        <DotCircle
            id={`dot-${payload.label}-${valueStringify(value)}-${idStringify(sectionId)}`}
            cx={cx}
            cy={cy}
            r={4}
            fill={color}
        />
    );
};

export const ClickableDot: FC<Props> = (props) => {
    const {cx, cy, fill, payload, value, originalColor, selectedDot, parentId, sectionId} = props;

    const isEcosystemRhs = useContext(IsEcosystemRhsContext);
    const fullUrl = useContext(FullUrlContext);
    const {url} = useRouteMatch();
    const {formatMessage} = useIntl();
    const {add: addToast} = useToaster();
    const [color, setColor] = useState(fill);

    const ecosystemQuery = isEcosystemRhs ? '' : buildQuery(parentId, sectionId);
    const [,, valueString] = dotStringify(cx, cy, value);

    const handleDotClick = (event: any) => {
        // const itemId = `_${cxString}_${cyString}_${payload.label}_${valueString}`;
        const itemId = `dot-${payload.label}-${valueString}-${idStringify(sectionId)}`;
        const name = `${payload.label}: ${value}`;
        const path = buildToForCopy(buildTo(fullUrl, itemId, ecosystemQuery, url));
        copyToClipboard(formatUrlAsMarkdown(path, name));
        addToast({content: formatMessage({defaultMessage: 'Copied!'})});
    };

    useEffect(() => {
        // if (selectedDot.x === cx && selectedDot.y === cy) {
        if (selectedDot.label === payload.label && valueStringify(selectedDot.value) === valueStringify(value)) {
            setColor('#F4B400');
        } else {
            setColor(originalColor);
        }
    }, [selectedDot]);

    return (
        <DotCircle
            cx={cx}
            cy={cy}
            r={7}
            fill={color}
            onClick={handleDotClick}
        />
    );
};

const DotCircle = styled.circle`
    cursor: pointer;    
`;