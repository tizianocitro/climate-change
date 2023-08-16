import React from 'react';
import {LoadingOutlined} from '@ant-design/icons';

export const InfoIcon = () => (
    <i className='icon fa fa-info-circle'/>
);

export const HiddenIcon = () => (
    <i
        id='hidden-icon'
        className='icon fa fa-ban'
    />
);

export const LoadingIcon = (
    <LoadingOutlined
        style={{fontSize: 24}}
        spin={true}
    />
);

// type Props = {
//     id: string;
//     style?: any;
// };

// export const FingerPointingIcon: FC<Props> = ({id, style}) => {
//     const color = '#F4B400';
//     return (
//         <FiArrowUp
//             id={id}
//             size={24}
//             style={style ? {
//                 color,
//                 ...style,
//             } : {
//                 color,
//             }}
//         />
//     );
// };
