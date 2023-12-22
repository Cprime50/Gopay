import { getMyself } from 'api/users';
import { getMyAccount } from 'api/accounts';
import { getMyCard } from 'api/cards';
import { getMyHistory } from 'api/history';

// Fetch data fro all the sections 

const getUserInitialData = async () => {
    try {
        const data = {
            user: await getMyself(),
            account: await getMyAccount(),
            card: await getMyCard(),
            history: await getMyHistory(),

        };

        return data;
    } catch (err) {
        throw new Error(err);
    }
};

export default getUserInitialData;